package scatter

import (
	"fmt"
	"time"
)

type Scatter struct{
	AppBootTime time.Time 			// app 启动时间
	ConfigTableMap map[string]interface{} 	// 配置
	baseTime time.Time
}

func NewScatter() *Scatter {
	info := new(Scatter)

	return info
}

// Init config and times
func (this *Scatter)Init() {

	// 记录app启动时间
	this.AppBootTime = time.Now()

	// 获取 2015 的时间，该时间按照 2006的格式产生
	tmptime, _ := time.Parse("2006-01-02 15:04:05", "2015-06-19 17:01:47.590")

	// app 启动的时间，时间格式为从 2015 年开始计算
	app_epoch_time := this.AppBootTime.Sub(tmptime)
	fmt.Printf("app epoch time: %f,%s ", app_epoch_time.Nanoseconds(), tmptime)

	// 显示可视时间
	viewable_datetime := fmt.Sprintf("%d%d%d%d%d%d", this.AppBootTime.Year(), this.AppBootTime.Month(), this.AppBootTime.Day(), this.AppBootTime.Hour(), this.AppBootTime.Minute(), this.AppBootTime.Second())
	fmt.Println(time.Now().Format("2006-01-02_15:04:05"), "   ", viewable_datetime)

	// 加载配置表
	this.ConfigTableMap = loadConfig()
}

type RoutineFunc func(goroutine_index int, user_data interface{})

// do func using N goroutines
// @param routine_func
// @param resulting scatter-result-filename
func (this *Scatter)Do(routine_func RoutineFunc, times int, user_data interface{}) (string, error){

	var N = times
	sem := make(chan string, N)
	for i := 0; i < N; i++ {
		go func(index int) {

			tnow := time.Now()
			starttimeint := (int)(tnow.Sub(this.AppBootTime).Seconds() * 1000000)

			///////////////////////////////////////
			routine_func(index, user_data)
			///////////////////////////////////////

			tmptime := time.Now()
			subdur := tmptime.Sub(tnow)
			tmpint := (int)(subdur.Seconds() * 1000)
			fmt.Println("index: ", index, "consumed: ", tmpint)

			sem <- (fmt.Sprintf("%d,%d", starttimeint, tmpint))
		}(i)
	}

	//	var max = 0
	outString := "["
	for m := 0; m < N; m++ {
		tmp := <-sem
		outString += fmt.Sprintf("[%s]", tmp)
		if m == N-1 {

		} else {
			outString += ","
		}
	}
	outString += "]"

	result_filename, err := writeFileWithData(this.ConfigTableMap["out_dir"].(string), "./config/output.html", outString, N)
	if err != nil {
		fmt.Println("failed to write file with scatter result data:", err.Error())
		return "", err
	}

	return result_filename, nil
}
