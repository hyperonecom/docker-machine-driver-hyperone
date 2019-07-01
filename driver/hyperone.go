package hyperone

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/docker/machine/libmachine/drivers"
	"github.com/docker/machine/libmachine/log"
	"github.com/docker/machine/libmachine/mcnflag"
	"github.com/docker/machine/libmachine/ssh"
	"github.com/docker/machine/libmachine/state"
	openapi "github.com/hyperonecom/h1-client-go"
)

var version = "devel"

const (
	defaultSSHUser  = "guru"
	defaultImage    = "debian"
	defaultType     = "a1.micro"
	defaultDiskName = "os-disk"
	defaultDiskType = "ssd"
	defaultDiskSize = 20
)

// Driver represents the docker driver interface
type Driver struct {
	*drivers.BaseDriver
	Token    string
	Project  string
	Image    string
	VMID     string
	Type     string
	DiskID   string
	DiskName string
	DiskType string
	DiskSize int
}

// GetCreateFlags registers the flags
func (d *Driver) GetCreateFlags() []mcnflag.Flag {
	return []mcnflag.Flag{
		mcnflag.StringFlag{
			EnvVar: "HYPERONE_ACCESS_TOKEN_SECRET",
			Name:   "hyperone-access-token-secret",
			Usage:  "HyperOne Access Token Secret",
		},
		mcnflag.StringFlag{
			EnvVar: "HYPERONE_PROJECT",
			Name:   "hyperone-project",
			Usage:  "HyperOne Project",
		},
		mcnflag.StringFlag{
			EnvVar: "HYPERONE_SSH_USER",
			Name:   "hyperone-ssh-user",
			Usage:  "SSH Username",
			Value:  defaultSSHUser,
		},
		mcnflag.StringFlag{
			EnvVar: "HYPERONE_IMAGE",
			Name:   "hyperone-image",
			Usage:  "HyperOne Image",
			Value:  defaultImage,
		},
		mcnflag.StringFlag{
			EnvVar: "HYPERONE_TYPE",
			Name:   "hyperone-type",
			Usage:  "HyperOne VM Type",
			Value:  defaultType,
		},
		mcnflag.StringFlag{
			EnvVar: "HYPERONE_DISK_NAME",
			Name:   "hyperone-disk-name",
			Usage:  "HyperOne VM OS Disk Name",
			Value:  defaultDiskName,
		},
		mcnflag.StringFlag{
			EnvVar: "HYPERONE_DISK_TYPE",
			Name:   "hyperone-disk-type",
			Usage:  "HyperOne VM OS Disk Type",
			Value:  defaultDiskType,
		},
		mcnflag.IntFlag{
			EnvVar: "HYPERONE_DISK_SIZE",
			Name:   "hyperone-disk-size",
			Usage:  "HyperOne VM OS Disk Size",
			Value:  defaultDiskSize,
		},
	}
}

// DriverName returns the name of the driver
func (d *Driver) DriverName() string {
	return "hyperone"
}

func (d *Driver) SetConfigFromFlags(flags drivers.DriverOptions) error {
	d.Token = flags.String("hyperone-access-token-secret")
	d.Project = flags.String("hyperone-project")
	d.Image = flags.String("hyperone-image")
	d.Type = flags.String("hyperone-type")
	d.DiskType = flags.String("hyperone-disk-type")
	d.DiskSize = flags.Int("hyperone-disk-size")
	d.DiskName = flags.String("hyperone-disk-name")
	d.SSHUser = flags.String("hyperone-ssh-user")

	if d.Token == "" {
		return fmt.Errorf("hyperone driver requires the --hyperone-access-token-secret option")
	}

	if d.Project == "" {
		return fmt.Errorf("hyperone driver requires the --hyperone-project option")
	}

	return nil
}

func (d *Driver) GetURL() (string, error) {
	ip, err := d.GetIP()
	if err != nil {
		return "", err
	}
	if ip == "" {
		return "", nil
	}
	return fmt.Sprintf("tcp://%s:2376", ip), nil
}

func (d *Driver) GetSSHHostname() (string, error) {
	return d.GetIP()
}

func (d *Driver) GetState() (state.State, error) {

	vm, _, err := d.getClient().VmApi.VmShow(context.TODO(), d.VMID)

	if err != nil {
		return state.Error, err
	}
	switch vm.State {
	case "Processing":
		return state.Starting, nil
	case "Running":
		return state.Running, nil
	case "Off":
		return state.Stopped, nil
	}
	return state.None, nil
}

func (d *Driver) Create() error {
	log.Infof("Creating HyperOne VM...")

	publicKey, err := d.createSSHKey()
	if err != nil {
		return err
	}

	client := d.getClient()

	options := openapi.VmCreate{
		Name:     d.MachineName,
		Service:  d.Type,
		Image:    d.Image,
		Username: d.SSHUser,
		SshKeys:  []string{strings.TrimSpace(publicKey)},
		Disk: []openapi.VmCreateDisk{
			{
				Name:    d.DiskName,
				Service: d.DiskType,
				Size:    float32(d.DiskSize),
			},
		},
	}

	vm, _, err := client.VmApi.VmCreate(context.TODO(), options)
	if err != nil {
		return err
	}

	hdds, _, err := client.VmApi.VmListHdd(context.TODO(), vm.Id)
	if err != nil {
		return err
	}

	d.VMID = vm.Id
	d.IPAddress = vm.Fqdn
	d.DiskID = hdds[0].Disk.Id

	return nil
}

func (d *Driver) Start() error {
	_, _, err := d.getClient().VmApi.VmActionStart(context.TODO(), d.VMID)
	return err
}

func (d *Driver) Stop() error {
	_, _, err := d.getClient().VmApi.VmActionStop(context.TODO(), d.VMID)
	return err
}

func (d *Driver) Restart() error {
	_, _, err := d.getClient().VmApi.VmActionRestart(context.TODO(), d.VMID)
	return err
}

func (d *Driver) Kill() error {
	_, _, err := d.getClient().VmApi.VmActionTurnoff(context.TODO(), d.VMID)
	return err
}

func (d *Driver) Remove() error {
	options := openapi.VmDelete{
		RemoveDisks: []string{d.DiskID},
	}

	if resp, err := d.getClient().VmApi.VmDelete(context.TODO(), d.VMID, options); err != nil {
		if resp.StatusCode == 404 {
			log.Infof("HyperOne VM doesn't exist, assuming it is already deleted")
		} else {
			return err
		}
	}

	return nil
}

// NewDriver returns a new driver
func NewDriver(hostName, storePath string) *Driver {
	return &Driver{
		BaseDriver: &drivers.BaseDriver{},
	}
}

func (d *Driver) getClient() *openapi.APIClient {
	cfg := openapi.NewConfiguration()
	cfg.AddDefaultHeader("x-project", d.Project)
	cfg.AddDefaultHeader("authorization", "Bearer "+d.Token)
	cfg.AddDefaultHeader("Prefer", "respond-async,wait=3600")

	cfg.UserAgent = fmt.Sprintf("docker-machine-driver-%s/%s", d.DriverName(), version)

	return openapi.NewAPIClient(cfg)
}

func (d *Driver) createSSHKey() (string, error) {
	if err := ssh.GenerateSSHKey(d.GetSSHKeyPath()); err != nil {
		return "", err
	}

	publicKey, err := ioutil.ReadFile(d.publicSSHKeyPath())
	if err != nil {
		return "", err
	}

	return string(publicKey), nil
}

func (d *Driver) publicSSHKeyPath() string {
	return d.GetSSHKeyPath() + ".pub"
}
