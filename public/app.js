var getQrCode = function () {
    $('#step1').addClass('hidden');
    $('#result').removeClass('hidden').addClass('hidden');
    
    user = $('#inpLogin').val();
    if (user.length == 0) return;
    

    $.post('/aj/genQRCode', { User: user })
        .done(function (data) {
            if (data.ok) {
                // show qrcode
                console.log(data.txt);
                $('#imgQRCode').attr('src', '/public/qrcodes/' + data.txt);
                $('#step2').removeClass('hidden');
                
            } else {
                console.log(data.txt);
            }
        })
        .fail(function (data) {
            console.log("ajax call failed: " + JSON.stringify(data));
        });
}

var testCode = function () {
    code = $('#inpCode').val();
    if (code.length == 0) {
        return
    }

    $.post('/aj/checkCode', { Code: code })
        .done(function (data) {
            if (data.ok) {
                if (data.txt == "true") {
                    $('#imgResult').attr('src','public/w.gif')                    
                } else {
                    $('#imgResult').attr('src','public/l.gif')
                }
                console.log(data.txt);
                 $('#result').removeClass('hidden');
            } else {
                console.log(data.txt);
            }
        })
        .fail(function (data) {
            console.log("ajax call failed: " + JSON.stringify(data));
        });



}



$(document).ready(function () {
    $('#btnLogin').click(getQrCode);
    $('#btnToStep3').click(function () {
        $('#step2').addClass('hidden');
        $('#step3').removeClass('hidden');
    });
    $('#btnTestCode').click(testCode);
})