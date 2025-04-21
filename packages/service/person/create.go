package person

import (
	"fedilist/packages/jsonld"
	"fedilist/packages/model/action"
	"fedilist/packages/model/list"
	"net/http"
)

func (ls PersonService) handleCreateList(ca action.Create, w http.ResponseWriter) {
	list, err := ls.listService.Create(func(ilv *list.ItemListValues) {
		ilv.Name = ca.Object().Name()
		ilv.Description = ca.Object().Description()
		ilv.Url = ca.Object().Url()
		ilv.Tags = ca.Object().Tags()
		ilv.Hooks = ca.Object().Hooks()
	})
	if err != nil {
		panic(err)
	}
	w.Write(jsonld.MarshalIndent(list))
}
