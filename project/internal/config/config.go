package config

import "log"
import "strconv"
import "time"
import "github.com/joho/godotenv"
import "os"

func SetConfig() (int, map[string]time.Duration) {
	//загрузка файла окружения
	err := godotenv.Load("../.env")
	if err != nil {
		log.Println(err)
	}

	//переменные окружения из файла .env	
	COMPUTING_POWER, errComputingPower 			 	:= os.LookupEnv("COMPUTING_POWER")
    TIME_ADDITION_MS, errAdditionTime  			 	:= os.LookupEnv("TIME_ADDITION_MS")
    TIME_SUBTRACTION_MS, errSubtractionTime      	:= os.LookupEnv("TIME_SUBTRACTION_MS")
    TIME_MULTIPLICATIONS_MS, errMultiplicationTime  := os.LookupEnv("TIME_MULTIPLICATIONS_MS")
    TIME_DIVISIONS_MS, errDivisionTime       		:= os.LookupEnv("TIME_DIVISIONS_MS")
    
	if !errComputingPower {
		log.Println("COMPUTING_POWER not in .env, default = 5")
		COMPUTING_POWER = "5"
	}
	if !errAdditionTime {
		log.Println("TIME_ADDITION_MS not in .env, default = 1")
		TIME_ADDITION_MS="1"
	}
	if !errSubtractionTime {
		log.Println("TIME_SUBTRACTION_MS not in .env, default = 1")
		TIME_SUBTRACTION_MS="1"
	}
	if !errMultiplicationTime {
		log.Println("TIME_MULTIPLICATIONS_MS not in .env, default = 1")
		TIME_MULTIPLICATIONS_MS="1"
	}
	if !errDivisionTime {
		log.Println("TIME_DIVISIONS_MS not in .env, default = 1")
		TIME_DIVISIONS_MS="1"
	}
    
	countExec, _ := strconv.Atoi(COMPUTING_POWER)
	
	additionTime, _         := strconv.Atoi(TIME_ADDITION_MS)
    subtractionTime, _      := strconv.Atoi(TIME_SUBTRACTION_MS)
    multiplicationTime, _   := strconv.Atoi(TIME_MULTIPLICATIONS_MS)
    divisionTime, _         := strconv.Atoi(TIME_DIVISIONS_MS)
	
    operationsTime := map[string]time.Duration {
        "+": time.Duration(additionTime),
        "-": time.Duration(subtractionTime),
        "*": time.Duration(multiplicationTime),
        "/": time.Duration(divisionTime),
    }
	
	return countExec, operationsTime
}