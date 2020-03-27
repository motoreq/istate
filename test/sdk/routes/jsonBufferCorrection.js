
var CorrectBuffer = function(obj) {
    Object.keys(obj).forEach(function(key) {
        var val = obj[key];
        if(val.type && val.type === "Buffer") {
            obj[key] = Buffer.from(val)
        } else if(Object.keys(val).length > 0) {
            obj[key] = CorrectBuffer(val)
        }
    });
    return obj
}

module.exports = CorrectBuffer