// Copyright 2023 The Ryan SU Authors (https://github.com/suyuan32). All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package casbin

import (
	"github.com/casbin/casbin/v2"
	"github.com/toutmost/admin-common/config"
	"log"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	entadapter "github.com/casbin/ent-adapter"
	rediswatcher "github.com/casbin/redis-watcher/v2"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
)

// CasbinConf is the configuration structure for Casbin
type CasbinConf struct {
	ModelText string `json:"ModelText,optional,env=CASBIN_MODEL_TEXT"`
}

// NewCasbin returns Casbin enforcer.
func (l CasbinConf) NewCasbin(dbType, dsn string) (*casbin.Enforcer, error) {
	adapter, err := entadapter.NewAdapter(dbType, dsn)
	logx.Must(err)

	var text string
	if l.ModelText == "" {
		text = `
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = r.sub == p.sub && keyMatch2(r.obj,p.obj) && r.act == p.act
		`
	} else {
		text = l.ModelText
	}

	m, err := model.NewModelFromString(text)
	logx.Must(err)

	enforcer, err := casbin.NewEnforcer(m, adapter)
	logx.Must(err)

	err = enforcer.LoadPolicy()
	logx.Must(err)

	return enforcer, nil
}

// MustNewCasbin returns Casbin enforcer. If there are errors, it will exist.
func (l CasbinConf) MustNewCasbin(dbType, dsn string) *casbin.Enforcer {
	csb, err := l.NewCasbin(dbType, dsn)
	if err != nil {
		logx.Errorw("initialize Casbin failed", logx.Field("detail", err.Error()))
		log.Fatalf("initialize Casbin failed, error: %s", err.Error())
		return nil
	}

	return csb
}

// MustNewRedisWatcher returns redis watcher. If there are errors, it will exist.
// f function will be called if the policies are updated.
func (l CasbinConf) MustNewRedisWatcher(c redis.RedisConf, f func(string2 string)) persist.Watcher {
	w, err := rediswatcher.NewWatcher(c.Host, rediswatcher.WatcherOptions{
		Options: redis2.Options{
			Network:  "tcp",
			Password: c.Pass,
		},
		Channel:    config.RedisCasbinChannel,
		IgnoreSelf: false,
	})
	logx.Must(err)

	err = w.SetUpdateCallback(f)
	logx.Must(err)

	return w
}

// MustNewCasbinWithRedisWatcher returns Casbin Enforcer with Redis watcher.
func (l CasbinConf) MustNewCasbinWithRedisWatcher(dbType, dsn string, c redis.RedisConf) *casbin.Enforcer {
	cbn := l.MustNewCasbin(dbType, dsn)
	w := l.MustNewRedisWatcher(c, func(data string) {
		rediswatcher.DefaultUpdateCallback(cbn)(data)
	})
	err := cbn.SetWatcher(w)
	logx.Must(err)
	err = cbn.SavePolicy()
	logx.Must(err)
	return cbn
}

// MustNewOriginalRedisWatcher returns redis watcher which uses original go redis. If there are errors, it will exist.
// f function will be called if the policies are updated.
func (l CasbinConf) MustNewOriginalRedisWatcher(c config.RedisConf, f func(string2 string)) persist.Watcher {
	w, err := rediswatcher.NewWatcher(c.Host, rediswatcher.WatcherOptions{
		Options: redis2.Options{
			Network:  "tcp",
			Username: c.Username,
			Password: c.Pass,
		},
		Channel:    "/casbin",
		IgnoreSelf: false,
	})
	logx.Must(err)

	err = w.SetUpdateCallback(f)
	logx.Must(err)

	return w
}

// MustNewCasbinWithOriginalRedisWatcher returns Casbin Enforcer with original Redis watcher.
func (l CasbinConf) MustNewCasbinWithOriginalRedisWatcher(dbType, dsn string, c config.RedisConf) *casbin.Enforcer {
	cbn := l.MustNewCasbin(dbType, dsn)
	w := l.MustNewOriginalRedisWatcher(c, func(data string) {
		rediswatcher.DefaultUpdateCallback(cbn)(data)
	})
	err := cbn.SetWatcher(w)
	logx.Must(err)
	err = cbn.SavePolicy()
	logx.Must(err)
	return cbn
}
