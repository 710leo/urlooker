package api

import (
	"github.com/710leo/urlooker/dataobj"
	"github.com/710leo/urlooker/modules/web/model"
)

func (this *Web) GetStrategies(req interface{}, resp *dataobj.StrategyResponse) error {
	strategies, err := model.GetAllStrategyByCron()
	if err != nil {
		resp.Message = err.Error()
	}
	stras := []*dataobj.Strategy{}
	for _, s := range strategies {
		stras = append(stras, &dataobj.Strategy{
			Id:          s.Id,
			Url:         s.Url,
			Enable:      s.Enable,
			IP:          s.IP,
			Keywords:    s.Keywords,
			Timeout:     s.Timeout,
			Creator:     s.Creator,
			ExpectCode:  s.ExpectCode,
			Note:        s.Note,
			Data:        s.Data,
			Tag:         s.Tag,
			MaxStep:     s.MaxStep,
			Times:       s.Times,
			Teams:       s.Teams,
			DingWebhook: s.DingWebhook,
		})
	}
	resp.Data = stras

	return nil
}
