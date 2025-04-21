package list

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fedilist/packages/jsonld"
	"fedilist/packages/model/action"
	"fedilist/packages/model/hook"
	"fedilist/packages/model/list"
	"fedilist/packages/service/cron"
	"fedilist/packages/util"
	"time"
)

func (s ListService) Create(fs ...func(*list.ItemListValues)) (list.ItemList, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	fs = append(fs, func(ilv *list.ItemListValues) {
		ilv.Key = base64.StdEncoding.EncodeToString(publicKey)
	})
	l, err := s.store.Insert(list.Create(fs...))
	if err != nil {
		return l, err
	}
	err = s.store.StoreKey(l, privateKey.Seed())
	if err != nil {
		return l, err
	}
	for _, h := range l.Hooks() {
		switch ch := h.(type) {
		case hook.CronHook:
			ea := action.CreateExecute(func(ev *action.ExecuteValues) {
				ev.Agent = l
				ev.StartTime = time.Now()
				ev.TargetRunner = ch.Runner()
				ev.RunnerAction = ch.RunnerAction()
				ev.RunnerActionConfig = ch.RunnerActionConfig()
			})
			pk, err := s.store.GetKey(l)
			if err != nil {
				panic(err)
			}
			b := jsonld.MarshalIndent(util.Sign(ea, pk))
			s.cronService.AddJob(cron.CronJob{
				Crontab: ch.CronTab(),
				Message: b,
			})
		}
	}
	return l, nil
}
