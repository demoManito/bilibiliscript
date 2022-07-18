## 同事吧盖楼

- 使用注意 ⚠️：
  - 使用的时候收敛一点，别太激进～
  - 设置 `X-CSRF`、`Cookie` 参数从请求头中获取
  - 设置 `Payload` 参数从中获取请求体
  - 以上参数均在 `config.yml` 中设置


- 功能说明：
  - **支持设置目标楼层**：
    - 需将 `target_floor.enable` 设置为 true 开启匹配楼层，默认为 false 表示关闭
    - 设置 `target_floor` 盖中目标楼层自动终止盖楼
    - 参数为数组，支持设置多个楼层
  - **支持规则匹配目标楼层**：
    - 需将 `target_floor_rule.enable` 设置为 true 开启规则匹配楼层，默认为 false 表示关闭
    - 设置 `target_floor_rule`，填写规则(`rule`) 和匹配目标数字(`target`) 即可
      1. 规则1: [rule: 1, target: 2] 表示匹配 2 的倍数
      2. 规则2: [rule: 2, target: 2] 表示匹配楼层数中包含 2 的楼层
  - **支持盖指定区间楼层**：
    - 需将 `target_floor_scope.enable` 设置为 true 开启盖指定区间楼层，默认为 false 表示关闭
    - 设置 `url` 评论列表接口，区间 [`min`, `max`]
  - **支持楼层触发器**：
    - 需将 `trigger_building.enable` 设置为 true 开启触发器，默认为 false 表示关闭
    - 设置 `num` 为触发的目标楼层，当前最大楼层数 >= `num` 则开始盖楼
    - 注意 ⚠️：此功能与**定时任务**互斥，同时设置默认只使用**楼层触发器**
  - **支持定时任务**：
    - 需将 `timing.enable` 设置为 true 开启定时任务，默认为 false 表示关闭
    - 设置 `timing.start_time`、`timing.end_time` 可定时开始和定时结束（可单独设置 `start_time` 或 `end_time` 两者无关联关系）
    - 日期格式：2022-07-22 23:00:00
    - 注意 ⚠️：此功能与**楼层触发器**互斥，同时设置默认只使用**楼层触发器**
  - **支持短信提醒**：
    - 需将 `sms_server.enable` 设置为 true 开启短信提醒服务，默认为 false 表示关闭
    - 需开通对应云厂商的短信服务，并填写对应参数
  - **支持多贴同时盖楼**：
  ```go
    func example() {
        builds := []*building.Building{
            building.New("./config1.yml"), 
            building.New("./config2.yml"),
        }
        building.RunBuilds(builds) // 开始同时盖楼
    }
  ```
  
- 如何启动：
  - 直接执行单测：`run_test.go`
  - 外部调用：`config.yml` 填写自己工程目录下文件路径的即可
  ```go
    // 单帖盖楼
    building.New("./config.example.yml").Run()
  
    // 多贴同时盖楼
    builds := []*building.Building{
        building.New("./config.example1.yml"), // NOTICE: config 不存在
        building.New("./config.example2.yml"), // NOTICE: config 不存在
    }
    building.RunBuilds(builds)
  ```