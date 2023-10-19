package file

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"
)

func AnalysisLogByLine(path string, detail bool, logStartTime, logEndTime string) error {
	f, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	ngLogStartTime := ""
	if logStartTime != "" {
		ngLogStartTime = time2ngTimeLocal(logStartTime)
	}
	ngLogEndTime := ""
	if logEndTime != "" {
		ngLogEndTime = time2ngTimeLocal(logEndTime)
	}

	if logStartTime != "" && logEndTime != "" {
		ok, err := endTimeGtStartTime(logStartTime, logEndTime)
		if err != nil {
			return err
		} else if !ok {
			return errors.New("日志筛选截止时间应大于开始时间")
		}
	}

	firstLogTime := ""
	lastLogTime := ""
	logTime := ""
	buf := bufio.NewScanner(f)
	num := 0
	numPerMin := make(map[string]uint)
	for {
		//当日志文件全部读完时，结束并打印日志统计信息
		if !buf.Scan() {
			lastLogTime = logTime
			maxKey, maxValue := findMaxValue(numPerMin, detail)
			fmt.Println("每分钟最大请求量: ", maxValue)
			fmt.Println("请求所处的时间点: ", maxKey)
			break
		}
		line := buf.Text() //获取每一行

		//在有筛序时间的情况下，没打到起始时间的日志全部忽略
		if ngLogStartTime != "" && num == 0 && !strings.Contains(line, ngLogStartTime) {
			if detail {
				fmt.Println("skip: ", line)
			}
			continue
		}

		logTime := strings.Split(line, " +0800")[0]
		if num == 0 {
			firstLogTime = logTime
		}
		//fmt.Println(logTime)
		logTimeTmp := strings.Split(logTime, ":")
		logTimeIgnoreSecond := logTimeTmp[0] + logTimeTmp[1] + logTimeTmp[2]
		//fmt.Println(logTimeIgnoreSecond)
		numPerMin[logTimeIgnoreSecond]++
		num++

		//达到截止时间的日志时，结束读取,并打印日志统计信息
		if ngLogEndTime != "" && strings.Contains(line, ngLogEndTime) {
			lastLogTime = logTime
			maxKey, maxValue := findMaxValue(numPerMin, detail)
			fmt.Println("每分钟最大请求量: ", maxValue)
			fmt.Println("请求所处的时间点: ", maxKey)
			break
		}
	}
	fmt.Println("日志开始时间:", firstLogTime)
	fmt.Println("日志结束时间:", lastLogTime)
	fmt.Println("nginx日志行数: ", num)
	return nil
}

func findMaxValue(m map[string]uint, detail bool) (string, uint) {
	var maxKey string
	var maxValue uint

	for k, v := range m {
		if detail {
			fmt.Println("时间: ", k, " 请求数: ", v)
		}

		if v > maxValue {
			maxValue = v
			maxKey = k
		}
	}

	return maxKey, maxValue
}

func time2ngTimeLocal(t string) string {
	mouthMap := make(map[string]string)
	mouthMap["01"] = "Jan"
	mouthMap["02"] = "Feb"
	mouthMap["03"] = "Mar"
	mouthMap["04"] = "Apr"
	mouthMap["05"] = "May"
	mouthMap["06"] = "Jun"
	mouthMap["07"] = "Jul"
	mouthMap["08"] = "Aug"
	mouthMap["09"] = "Sep"
	mouthMap["10"] = "Oct"
	mouthMap["11"] = "Nov"
	mouthMap["12"] = "Dec"
	// 2023-10-19-12:12 ---> 17/Oct/2023:14:35
	timeTmp := strings.Split(t, "-")
	return fmt.Sprintf("%s/%s/%s:%s", timeTmp[2], mouthMap[timeTmp[1]], timeTmp[0], timeTmp[3])
}

func endTimeGtStartTime(startTime, endTime string) (bool, error) {
	st, err := time.Parse("2006-01-02-15:04", startTime)
	if err != nil {
		return false, err
	}

	et, err := time.Parse("2006-01-02-15:04", endTime)
	if err != nil {
		return false, err
	}
	return et.Unix() >= st.Unix(), nil
}
