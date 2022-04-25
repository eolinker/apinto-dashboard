
function PropertiesToArray(data){

    let tName =data["type"]
    if (tName !== "object"){
        throw {toString: function() { return "root only allow object"; } };
    }
    if (data["additionalProperties"] ){
        throw {toString: function() { return "root not allow additionalProperties"; } };
    }
    return readObject(data)
}

function readProperties(parent,properties){
    let list = []
    for (const name in properties) {
        list.push(toItem(parent,properties,name))
    }
    return list
}

function readArrayItems(obj){
    switch (obj["type"]){
        case "object":{
            return readObject(obj)
        }
    }
    return obj
}
function readObject(data){
    var value = {...data}
    delete value["additionalProperties"]
    delete value["required"]
    delete value["items"]
    if (data["additionalProperties"]){
        value["type"] = "map"
        value["items"] = readItems(data["additionalProperties"])
        delete value["additionalProperties"]
    }else{
        value["type"] = "object"
        value["properties"] = readProperties(data,data["properties"])
    }
    return value
}
function readItems(data){
    if (data["type"] === "object"){
        return readObject(data)
    }else{
        return data
    }
}
function toItem(parent,properties,name){

    var data = properties[name]

    var value = {"name":name,...data}

    switch (data["type"]){
        case "object":{
            let d = readObject(value)
            value = {
                name:name,
                ...d
            }
            break
        }
        case "array":{
            value["type"] = "array"
            value["items"] = readArrayItems(data["items"])
            delete value["additionalProperties"]
            delete value["required"]
            delete value["properties"]

            break
        }
        default:{
            delete value["additionalProperties"]
            delete value["required"]
            delete value["items"]
            delete value["properties"]
            break
        }
    }
    if (parent["required"] && parent["required"].find(element => element === name)){
        value["required"] = true
    }

    return value
}

// let ddddd = PropertiesToArray(demoSchema)
// console.log(ddddd)