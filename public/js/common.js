const $ = mdui.$

function setUserInfo(userInfo) {
    userInfo['auth'] = userInfo['token'] + "@@" + userInfo['uid']
    localStorage.setItem("userInfo", JSON.stringify(userInfo))
}

function getUserInfo() {
    let userInfo = JSON.parse(localStorage.getItem("userInfo"))
    if (userInfo === undefined || userInfo === null) {
        window.location.href = "/view/user/login"
        return null
    }
    let uid = userInfo['uid']
    if (uid !== null && uid !== undefined && uid !== '') {
        return userInfo
    }
    return null
}

function getAuthToken() {
    let userInfo = getUserInfo()
    let auth = ""
    if (userInfo !== null) {
        userId = userInfo.uid
        token = userInfo.token
    }
    auth = token + "@@" + userId
    return auth
}

function cleanUserInfo() {
    localStorage.removeItem("userInfo")
}

function isEmail(email) {
    return String(email)
        .toLowerCase()
        .match(
            /^(([^<>()[\]\\.,:\s@"]+(\.[^<>()[\]\\.,:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/
        )
}

function isMobile(mobile) {
    let reg = /^1[3-9]\d{9}$/
    return reg.test(mobile)
}

function parseDescriptionStringToHtml() {
    $('#feed-item-list').find('.feed-item-description-tag').each(function (index, element) {
        if (!$(this).hasClass('parsed')) {
            $(this).html($(this).text())
            $(this).addClass('parsed')
        }
    })
}

function getSubChannelListTmpl() {
    let userInfo = getUserInfo()
    let userId = ""
    if (userInfo !== null) {
        userId = userInfo.uid
    }
    $.ajax({
        method: 'GET',
        headers: {
            Authorization: getAuthToken()
        },
        url: '/view/feed/sub_list',
        data: {
            userId: userId,
            start: 0,
            size: 10,
        },
        success: function (data) {
            $('#sub-channel-drawer-list').append(data)
        }
    });
}

function markedItem(id) {
    let userInfo = getUserInfo()
    let userId = ""
    if (userInfo !== null) {
        userId = userInfo.uid
    }
    let postData = {
        userId: userId,
        itemId: id,
    }
    $.ajax({
        method: 'POST',
        headers: {
            Authorization: getAuthToken()
        },
        url: '/v1/api/feed/item/mark',
        data: JSON.stringify(postData),
        success: function (data) {
            let jsonData = JSON.parse(data)
            if (jsonData.error === 0) {
                if (jsonData.data[0] === 1) {
                    $('#marked-icon-' + id).addClass('mdui-text-color-theme')
                } else {
                    $('#marked-icon-' + id).removeClass('mdui-text-color-theme')
                }
            }
            mdui.snackbar({
                message: jsonData.msg,
                position: 'top',
            })
        }
    });
}

function share(id) {
    let domain = window.location.host
    let shareUrl = 'http://' + domain + '/view/f/i/s/' + id
    $('#clipboard-temp-holder-' + id).text(shareUrl)
    copyToClickBoard('clipboard-temp-holder-' + id)

}

function copyToClickBoard(elemId) {
    var content = document.getElementById(elemId).innerHTML;
    navigator.clipboard.writeText(content)
        .then(() => {
            mdui.snackbar({
                message: '已复制到剪贴板',
                position: 'top'
            })
        })
        .catch(err => {
            console.log('Something went wrong', err);
        })
}

function optimizeImg() {
    const img = document.querySelectorAll("img")
    img.forEach(item => {
        item.style.maxWidth = "100%"
        item.style.height = "auto"
    })
    const divs = document.querySelectorAll("div")
    divs.forEach(item => {
        if (item.style.width.replace("px", "") > window.innerWidth) {
            item.style.maxWidth = "100%"
            item.style.height = "auto"
        }
    })
}