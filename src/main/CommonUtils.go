package main

import "time"

func truncateTimeToDate(from time.Time)(time.Time){
	to := from.Truncate( time.Hour ).Add( - time.Duration(from.Hour()) * time.Hour )
	return to
}