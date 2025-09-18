# API Guide 

## Default format for Protocol Messages:

### REQUEST 
```json
{
    "method": "method_name",
    "timestamp": "timestamp_value",
    "data": {
        // method-specific parameters
    }
}
```

### RESPONSE
```json
{
    "method": "method_name",
    "status": "success" | "error",
    "message": "descriptive_message",
    "timestamp": "timestamp_value",
    "data": {
        // method-specific response data
    }
}
```
