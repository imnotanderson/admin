function test1() {
    sendAjaxRequest("/req?q=test&p=anderson",function (s) {
        return s;
    })
}