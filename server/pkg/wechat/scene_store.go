package wechat

import (
	"sync"
	"time"
)

type sceneStatus string

const (
	ScenePending   sceneStatus = "pending"
	SceneConfirmed sceneStatus = "confirmed"
)

type sceneInfo struct {
	userID    uint
	status    sceneStatus
	expiresAt time.Time
}

var (
	scenes = &sync.Map{}
)

// CreateScene 创建扫码登录场景
func CreateScene(scene string) {
	scenes.Store(scene, &sceneInfo{
		status:    ScenePending,
		expiresAt: time.Now().Add(5 * time.Minute),
	})
}

// GetSceneStatus 查询场景状态，返回 (status, userID, exists)
func GetSceneStatus(scene string) (string, uint, bool) {
	val, ok := scenes.Load(scene)
	if !ok {
		return "", 0, false
	}
	info := val.(*sceneInfo)
	if time.Now().After(info.expiresAt) {
		scenes.Delete(scene)
		return "", 0, false
	}
	return string(info.status), info.userID, true
}

// ConfirmScene 小程序确认授权
func ConfirmScene(scene string, userID uint) bool {
	val, ok := scenes.Load(scene)
	if !ok {
		return false
	}
	info := val.(*sceneInfo)
	if time.Now().After(info.expiresAt) || info.status != ScenePending {
		return false
	}
	info.userID = userID
	info.status = SceneConfirmed
	return true
}
