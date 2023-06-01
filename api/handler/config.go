package handler

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/jxo-me/netx/x/config"
	"os"
)

var (
	Config = hConfig{}
)

type hConfig struct{}

type GetConfigReq struct {
	g.Meta `path:"/" method:"get" tags:"Config" summary:"Get current config."`
	// output format, one of yaml|json, default is json.
	// in: query
	Format string `form:"format" json:"format" description:"output format, one of yaml|json, default is yaml."`
}

type GetConfigRes struct {
	Config *config.Config
}

func (h *hConfig) GetConfig(ctx context.Context, req *GetConfigReq) (res *GetConfigRes, err error) {
	res = &GetConfigRes{}
	res.Config = config.Global()
	glog.Info(ctx, "GetConfig:", res)
	return res, nil
}

type SaveConfigReq struct {
	g.Meta `path:"/" method:"post" tags:"Config" summary:"Save current config to file (gost.yaml or gost.json)."`
	// output format, one of yaml|json, default is yaml.
	// in: query
	Format string `form:"format" json:"format" description:"output format, one of yaml|json, default is yaml."`
}

func (h *hConfig) SaveConfig(ctx context.Context, req *SaveConfigReq) (res *NullStructRes, err error) {
	res = &NullStructRes{}
	file := "gost.yaml"
	switch req.Format {
	case "json":
		file = "gost.json"
	default:
		req.Format = "yaml"
	}

	f, err := os.Create(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if err := config.Global().Write(f, req.Format); err != nil {
		return nil, err
	}

	return res, nil
}
