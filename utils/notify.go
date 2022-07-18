package utils

import (
	"time"

	"github.com/progrium/macdriver/cocoa"
	"github.com/progrium/macdriver/core"
	"github.com/progrium/macdriver/objc"
)

var (
	NSUserNotification_       = objc.Get("NSUserNotification")
	NSUserNotificationCenter_ = objc.Get("NSUserNotificationCenter")
)

type NSUserNotification struct {
	objc.Object
}

type NSUserNotificationCenter struct {
	objc.Object
}

type notifyController struct {
	app cocoa.NSApplication
}

func Notify(msg string) {
	notify := new(notifyController)
	notify.app = cocoa.NSApp_WithDidLaunch(func(_ objc.Object) {
		notification := NSUserNotification{NSUserNotification_.Alloc().Init()}
		notification.Set("title:", core.String("同事吧盖楼提醒 ⏰"))
		notification.Set("informativeText:", core.String(msg))
		center := NSUserNotificationCenter{NSUserNotificationCenter_.Send("defaultUserNotificationCenter")}
		center.Send("deliverNotification:", notification)
		notification.Release()
	})

	nsbundle := cocoa.NSBundle_Main().Class()
	nsbundle.AddMethod("__bundleIdentifier", func(_ objc.Object) objc.Object {
		return core.String("com.example.fake")
	})
	nsbundle.Swizzle("bundleIdentifier", "__bundleIdentifier")
	notify.app.SetActivationPolicy(cocoa.NSApplicationActivationPolicyRegular)
	notify.app.ActivateIgnoringOtherApps(true)
	// notify.hook()
	notify.app.Run()
}

func (nc *notifyController) hook() {
	go func() {
		timer := time.NewTimer(2 * time.Second)
		<-timer.C
		for {
			if nc.app.IsRunning() {
				nc.app.Terminate()
			}
		}
	}()
}
