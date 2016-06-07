var getQrCode = function () {
    user = $('#inpLogin').val();
    if (user.length == 0 ) return;
    
    $.post('/aj/genQRCode', { user: user})
        .done(function (data) {
            if (data.ok) {
                // show code
                console.log(data.txt);
                $('#imgQRCode').attr('src', '/public/qrcodes/'+data.txt);
                $('#step2').removeClass('hidden');                
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
})