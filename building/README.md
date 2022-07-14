## 同事吧盖楼

- 使用注意 ⚠️：
  - 使用的时候收敛一点，别太激进～
  - 设置 `X-CSRF`、`Cookie` 参数从请求头中获取
  - 设置 `Payload` 参数从中获取请求体
  - 以上参数均在 `config.yml` 中设置


- 功能说明：
  - 支持设置目标楼层：设置参数 `target_floor`，盖中目标楼层自动终止盖楼
  - 支持规则匹配目标楼层：设置参数 `target_floor_rule`，填写规则(`rule`) 和匹配目标数字(`target`) 即可
  - 支持设置定时任务：设置参数 `timing_start_time`、`timing_end_time`，可定时开始定时结束
  - 支持多贴同时盖楼：
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
  - 外部调用：(`config.yml` 填写自己工程目录下文件路径的即可)
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