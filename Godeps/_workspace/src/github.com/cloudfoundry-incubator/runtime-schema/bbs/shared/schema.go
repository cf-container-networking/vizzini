package shared

import (
	"path"
	"strconv"

	"github.com/cloudfoundry-incubator/runtime-schema/models"
)

const SchemaRoot = "/v1/"
const ExecutorSchemaRoot = SchemaRoot + "executor"
const FileServerSchemaRoot = SchemaRoot + "file_server"
const LRPStartAuctionSchemaRoot = SchemaRoot + "start"
const LRPStopAuctionSchemaRoot = SchemaRoot + "stop"
const StopLRPInstanceSchemaRoot = SchemaRoot + "stop-instance"
const ActualLRPSchemaRoot = SchemaRoot + "actual"
const DesiredLRPSchemaRoot = SchemaRoot + "desired"
const TaskSchemaRoot = SchemaRoot + "task"
const LockSchemaRoot = SchemaRoot + "locks"

func ExecutorSchemaPath(executorID string) string {
	return path.Join(ExecutorSchemaRoot, executorID)
}

func FileServerSchemaPath(segments ...string) string {
	return path.Join(append([]string{FileServerSchemaRoot}, segments...)...)
}

func LRPStartAuctionSchemaPath(lrp models.LRPStartAuction) string {
	return path.Join(LRPStartAuctionSchemaRoot, lrp.ProcessGuid, strconv.Itoa(lrp.Index))
}

func LRPStopAuctionSchemaPath(lrp models.LRPStopAuction) string {
	return path.Join(LRPStopAuctionSchemaRoot, lrp.ProcessGuid, strconv.Itoa(lrp.Index))
}

func StopLRPInstanceSchemaPath(stopInstance models.StopLRPInstance) string {
	return path.Join(StopLRPInstanceSchemaRoot, stopInstance.InstanceGuid)
}

func ActualLRPSchemaPathFromStopLRPInstance(stopInstance models.StopLRPInstance) string {
	return path.Join(ActualLRPSchemaRoot, stopInstance.ProcessGuid, strconv.Itoa(stopInstance.Index), stopInstance.InstanceGuid)
}

func ActualLRPSchemaPath(lrp models.ActualLRP) string {
	return path.Join(ActualLRPSchemaRoot, lrp.ProcessGuid, strconv.Itoa(lrp.Index), lrp.InstanceGuid)
}

func DesiredLRPSchemaPath(lrp models.DesiredLRP) string {
	return DesiredLRPSchemaPathByProcessGuid(lrp.ProcessGuid)
}

func DesiredLRPSchemaPathByProcessGuid(processGuid string) string {
	return path.Join(DesiredLRPSchemaRoot, processGuid)
}

func TaskSchemaPath(taskGuid string) string {
	return path.Join(TaskSchemaRoot, taskGuid)
}

func LockSchemaPath(lockName string) string {
	return path.Join(LockSchemaRoot, lockName)
}
