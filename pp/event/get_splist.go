package event

import (
	"context"
	"github.com/stratosnet/sds/pp/setting"
	"time"

	"github.com/stratosnet/sds/framework/core"
	"github.com/stratosnet/sds/msg/header"
	"github.com/stratosnet/sds/msg/protos"
	"github.com/stratosnet/sds/utils"

	"github.com/alex023/clock"
)

// GetPPList P node get PPList
func GetSPList() {
	utils.DebugLog("SendMessage(client.SPConn, req, header.ReqGetSPList)")
	SendMessageToSPServer(reqGetSPlistData(), header.ReqGetSPList)
}

// RspGetPPList
func RspGetSPList(ctx context.Context, conn core.WriteCloser) {
	utils.Log("get GetSPList RSP")
	var target protos.RspGetSPList
	if !unmarshalData(ctx, &target) {
		return
	}
	utils.Log("get GetSPList RSP", target.SpList)
	if target.Result.State != protos.ResultState_RES_SUCCESS {
		reloadSPlist()
		return
	}

	changed := false
	for _, sp := range target.SpList {
		existing, ok := setting.SPMap.Load(sp.P2PAddress)
		if ok {
			existingSp := existing.(setting.SPBaseInfo)
			if sp.P2PPubKey != existingSp.P2PPublicKey || sp.NetworkAddress != existingSp.NetworkAddress {
				changed = true
			}
		} else {
			changed = true
		}

		setting.SPMap.Store(sp.P2PAddress, setting.SPBaseInfo{
			P2PAddress:     sp.P2PAddress,
			P2PPublicKey:   sp.P2PPubKey,
			NetworkAddress: sp.NetworkAddress,
		})
	}
	if changed {
		setting.SPMap.Delete("unknown")
		setting.Config.SPList = nil
		setting.SPMap.Range(func(k, v interface{}) bool {
			sp := v.(setting.SPBaseInfo)
			setting.Config.SPList = append(setting.Config.SPList, sp)
			return true
		})
		if err := utils.WriteConfig(setting.Config, setting.ConfigPath); err != nil {
			utils.ErrorLog("Couldn't write config with updated SP list to yaml file", err)
		}
	}
}

func reloadSPlist() {
	utils.DebugLog("failed to get SPlist. retry after 3 second")
	newClock := clock.NewClock()
	newClock.AddJobRepeat(time.Second*3, 1, GetSPList)
}
