package api

import (
	cmodel "github.com/open-falcon/common/model"

	"github.com/open-falcon/graph/index"
)

// 删除一条索引，及其对应的graph索引缓存
func (this *Graph) DeleteIndex(counter *cmodel.GraphCounter, resp *cmodel.SimpleRpcResponse) error {
	err := index.DeleteCounter(counter.Endpoint, counter.Counter)
	if err != nil {
		resp.Code = 1
		resp.Msg = err.Error()
	}

	return nil
}
