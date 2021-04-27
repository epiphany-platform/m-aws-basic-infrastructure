package cmd

import (
	"errors"
	awsbi "github.com/epiphany-platform/e-structures/awsbi/v0"
	st "github.com/epiphany-platform/e-structures/state/v0"
	"github.com/epiphany-platform/e-structures/utils/load"
	"github.com/epiphany-platform/e-structures/utils/to"
	"io/ioutil"
	"os"
)

func ensureDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func checkAndLoad(stateFilePath string, configFilePath string) (*awsbi.Config, *st.State, error) {
	logger.Debug().Msgf("checkAndLoad(%s, %s)", stateFilePath, configFilePath)
	if _, err := os.Stat(stateFilePath); os.IsNotExist(err) {
		return nil, nil, errors.New("state file does not exist, please run init first")
	}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return nil, nil, errors.New("config file does not exist, please run init first")
	}

	state, err := load.State(stateFilePath)
	if err != nil {
		return nil, nil, err
	}

	config, err := load.AwsBIConfig(configFilePath)
	if err != nil {
		return nil, nil, err
	}

	return config, state, nil
}

func backupFile(path string) error {
	logger.Debug().Msgf("backupFile(%s)", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	} else {
		backupPath := path + ".backup"

		input, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		err = ioutil.WriteFile(backupPath, input, 0644)
		if err != nil {
			return err
		}
		return nil
	}
}

func produceOutput(m map[string]interface{}) *awsbi.Output {
	logger.Debug().Msgf("Received output map: %#v", m)

	// two internal intermediate data structures to hold extracted map values
	type intermediateDataDisk struct {
		id   string
		name string
		size int
	}
	type intermediateDataDiskAttachment struct {
		deviceName string
		dataDiskId string
		instanceId string
	}

	output := &awsbi.Output{
		VpcId: to.StrPtr(m["vpc_id"].(string)),
	}
	if prt, ok := m["private_route_table"]; ok {
		output.PrivateRouteTable = to.StrPtr(prt.(string))
	}
	for _, k := range m["private_subnet_ids"].([]interface{}) {
		output.PrivateSubnetIds = append(output.PrivateSubnetIds, k.(string))
	}
	for _, k := range m["public_subnet_ids"].([]interface{}) {
		output.PublicSubnetIds = append(output.PublicSubnetIds, k.(string))
	}
	for _, i := range m["vm_group"].([]interface{}) {
		vmGroup := i.(map[string]interface{})
		outputVmGroup := awsbi.OutputVmGroup{
			Name: to.StrPtr(vmGroup["vm_group_name"].(string)),
		}
		intermediateDataDisks := make([]intermediateDataDisk, 0)
		for _, j := range vmGroup["data_disks"].([]interface{}) {
			tempDataDisk := j.(map[string]interface{})
			intermediateDataDisks = append(intermediateDataDisks,
				intermediateDataDisk{
					id:   tempDataDisk["id"].(string),
					name: tempDataDisk["name"].(string),
					size: int(tempDataDisk["size"].(float64)),
				})
		}
		logger.Debug().Msgf("Intermediate data disks struct list: %#v", intermediateDataDisks)
		intermediateDataDiskAttachments := make([]intermediateDataDiskAttachment, 0)
		for _, j := range vmGroup["dd_attachments"].([]interface{}) {
			tempDataDiskAttachment := j.(map[string]interface{})
			intermediateDataDiskAttachments = append(intermediateDataDiskAttachments,
				intermediateDataDiskAttachment{
					deviceName: tempDataDiskAttachment["device_name"].(string),
					dataDiskId: tempDataDiskAttachment["data_disk_id"].(string),
					instanceId: tempDataDiskAttachment["instance_id"].(string),
				})
		}
		logger.Debug().Msgf("Intermediate data disk attachments struct list: %#v", intermediateDataDiskAttachments)
		for _, j := range vmGroup["vms"].([]interface{}) {
			tempVm := j.(map[string]interface{})
			outputVm := awsbi.OutputVm{
				Name:      to.StrPtr(tempVm["vm_name"].(string)),
				PublicIp:  to.StrPtr(tempVm["public_ip"].(string)),
				PrivateIp: to.StrPtr(tempVm["private_ip"].(string)),
			}
			vmId := tempVm["id"].(string)
			for _, dda := range intermediateDataDiskAttachments {
				if dda.instanceId == vmId {
					for _, dd := range intermediateDataDisks {
						if dd.id == dda.dataDiskId {
							outputVm.DataDisks = append(outputVm.DataDisks, awsbi.OutputDataDisk{
								Size:       to.IntPtr(dd.size),
								DeviceName: to.StrPtr(dda.deviceName),
							})
						}
					}
				}
			}
			outputVmGroup.Vms = append(outputVmGroup.Vms, outputVm)
		}
		output.VmGroups = append(output.VmGroups, outputVmGroup)
	}
	return output
}
