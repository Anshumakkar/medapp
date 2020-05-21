package db

// snippet-start:[dynamodb.go.create_item.imports]
import (
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-sdk-go/aws/credentials"
    "fmt"
)
// snippet-end:[dynamodb.go.create_item.imports]

// snippet-start:[dynamodb.go.create_item.struct]
// Create struct to hold info about new item
type Item struct {
    Year   int
    Title  string
    Plot   string
    Rating float64
}


type LoginData struct {
        Email       string
        Gender      string
        Password    string
        Username    string
        PhoneNumber string
	ID string
}

var svc *dynamodb.DynamoDB

// snippet-end:[dynamodb.go.create_item.struct]
func CreateClient() error {
    // snippet-start:[dynamodb.go.create_item.session]
    // Initialize a session that the SDK will use to load
    // credentials from the shared credentials file ~/.aws/credentials
    // and region from the shared configuration file ~/.aws/config.

creds := credentials.NewStaticCredentials("AKIAJWLDEWOE33G3TRCA","M6Uuec7klcMjIwHue9l12XWbya+UkSS8D/LVgOVv","")
sess, err := session.NewSession(&aws.Config{Credentials: creds,Region: aws.String("ap-south-1")})
if err!=nil{
	return err
}
    // Create DynamoDB client
    svc = dynamodb.New(sess)
    // snippet-end:[dynamodb.go.create_item.session]
    return nil
}

func GetItem(phoneNumber string) (*LoginData,error){
    item := LoginData{}

    tableName := "usersInfo"


    result, err := svc.GetItem(&dynamodb.GetItemInput{
        TableName: aws.String(tableName),
        Key: map[string]*dynamodb.AttributeValue{
            "PhoneNumber": {
                S: aws.String(phoneNumber),
            },
	    //"Email":{
	//	    S: aws.String("test@testing.com"),
	  //  },

        },
    })
    if err != nil {
        fmt.Println("Error got " ,err.Error())
        return nil,err
    }


    err = dynamodbattribute.UnmarshalMap(result.Item, &item)
    if err != nil {
		fmt.Println("Failed to unmarshal Record: ", err.Error())
        return nil,err
    }
    return &item,nil

}

func StoreItem(phoneNumber,email,username,password,gender,id string) error{
    // snippet-start:[dynamodb.go.create_item.assign_struct]
    item := LoginData{
        Email:   email,
        Gender:  gender,
        Password:   password,
        Username: username,
	PhoneNumber: phoneNumber,
	ID:id,
    }

    av, err := dynamodbattribute.MarshalMap(item)
    if err != nil {
        fmt.Println("Got error marshalling new user registration info :")
        fmt.Println(err.Error())
        return err
    }
    // snippet-end:[dynamodb.go.create_item.assign_struct]

    // snippet-start:[dynamodb.go.create_item.call]
    // Create item in table Movies
    tableName := "usersInfo"

    input := &dynamodb.PutItemInput{
        Item:      av,
        TableName: aws.String(tableName),
    }

    _, err = svc.PutItem(input)
    if err != nil {
        fmt.Println("Got error calling PutItem:")
        fmt.Println(err.Error())
        return err
    }


    fmt.Println("Successfully added '" + item.Username + "' (" + item.PhoneNumber + ") to table " + tableName)
    return nil
    // snippet-end:[dynamodb.go.create_item.call]
}
