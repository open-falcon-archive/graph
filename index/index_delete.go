package index

import (
	"fmt"
	"log"

	"github.com/open-falcon/common/utils"

	"github.com/open-falcon/graph/g"
)

func DeleteCounter(endpoint, counter string) error {
	if len(endpoint) < 1 || len(counter) < 1 {
		return fmt.Errorf("bad args, endpoint counter: %s, %s", endpoint, counter)
	}

	config := g.Config()
	// get mysqlConn
	conn, err := g.GetDbConn("DeleteIndex")
	if err != nil {
		if config.Debug {
			log.Println("[ERROR] get mysql conn fail", err)
		}
		return err
	}

	// get endpoint_id
	var endpointId int64 = -1
	var targetEndpoint string = ""
	err = conn.QueryRow("SELECT id,endpoint FROM endpoint WHERE endpoint = ?", endpoint).Scan(&endpointId, &targetEndpoint)
	if err != nil {
		if config.Debug {
			log.Println("[ERROR] query endpoint fail", err)
		}
		return err
	}
	if targetEndpoint != endpoint { // 防止注入
		return fmt.Errorf("bad arg: endpoint %s", endpoint)
	}

	// delete counter
	var endpointCounterId int64 = -1
	var targetCounter string = ""
	err = conn.QueryRow("SELECT id,counter FROM endpoint_counter WHERE endpoint_id = ? and counter = ?",
		endpointId, counter).Scan(&endpointCounterId, &targetCounter)
	if err != nil {
		if config.Debug {
			log.Println("[ERROR] query counter fail", err)
		}
		return err
	}
	if targetCounter != counter { // 防止注入
		return fmt.Errorf("bad arg: counter %s", counter)
	}

	_, err = conn.Exec("DELETE FROM endpoint_counter WHERE id = ?", endpointCounterId)
	if err != nil {
		if config.Debug {
			log.Println("[ERROR] delete counter fail", err)
		}
		return err
	}
	log.Printf("delete counter: %s/%s", endpoint, counter)

	// delete indexed cache
	key := utils.ChecksumOfPK2(endpoint, counter)
	indexedItemCache.Remove(key)

	// return
	return nil
}
