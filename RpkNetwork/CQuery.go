package RpkNetwork

import (
	"fmt"
	
)


func Query(data CData) CData {
	var re CData


	iAction:=data.Action

	switch  {
	case iAction >=API_REQUSET && iAction <9900:
			fmt.Println("api")
			return re
			
	case iAction ==ADD_USER:
			fmt.Println("api")
			return re			


	}

	return re
}

