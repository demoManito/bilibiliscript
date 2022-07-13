## 同事吧盖楼

- 使用注意 ⚠️：
  - 使用的时候收敛一点，别太猖狂
  - 设置 `X-CSRF`、`Cookie` 参数从请求头中获取
  - 设置 `Payload` 参数从中获取请求体
  - 以上参数均在 `config.yml` 中设置


- 功能说明：
  - 支持设置目标楼层：设置参数 `target_floor`，盖中目标楼层自动终止盖楼
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
    