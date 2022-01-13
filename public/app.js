new Vue({//tao mot doi tuong vueJS
    el: '#app',
    // yeu cau vue hien thi ung dung ben trong phan tu DOM bang id '#app'
    // doi tuong data la noi dat du lieu su dung trong ung dung
    data: {
        ws: null, // tao bien ws de luu tru Websoket
        newMsg: '', //Giu tin nhan moi de gui den server
        chatContent: '', // danh sach tin nhan dang chay tren man hinh
        email: null, // dia chi email de lay avatar
        username: null,
        joined: false // neu nhu email va username da dien thi se la true
    },

    created: function() {// ham dinh nghia de xu ly
        var self = this;// khai bao bien dang con tro this
        this.ws = new WebSocket('ws://' + window.location.host + '/ws');//tao mot ket noi Websocket den may chu va
        // mot trinh xu ly khi cac tin nhan moi duoc gui tu may chu
        this.ws.addEventListener('message', function(e) { //ham xu ly cac tin nhan den
            var msg = JSON.parse(e.data);
            self.chatContent += '<div class="chip">'
                + '<img src="' + self.gravatarURL(msg.email) + '">' // lay ra avata dua vao email
                + msg.username
                + '</div>'
                + emojione.toImage(msg.message) + '<br/>'; // Parse emojis

            var element = document.getElementById('chat-messages');
            element.scrollTop = element.scrollHeight; // tu dong day xuong duoi
        });
    },

    methods: { // dinh nghia cac chuc nang su dung trong VueJs
        send: function () { // xu ly gui tin nhan den server
            if (this.newMsg != '') { //ktra tin nhan trong
                this.ws.send( // dinh dang thong bao duoi dang doi tuong
                    JSON.stringify({ //sd stringify de may chu phan tich cu phap
                            email: this.email,
                            username: this.username,
                            message: $('<p>').html(this.newMsg).text() // loai bo html
                        }
                    ));
                this.newMsg = ''; // Reset newMsg
            }
        },

        join: function () { // ham bat buoc nguoi dung phai nhap email va username truoc khi gui tn
            if (!this.email) {
                Materialize.toast('You must enter an email', 2000);
                return
            }
            if (!this.username) {
                Materialize.toast('You must choose a username', 2000);
                return
            }
            this.email = $('<p>').html(this.email).text();
            this.username = $('<p>').html(this.username).text();
            this.joined = true;
        },

        gravatarURL: function(email) {
            return 'http://www.gravatar.com/avatar/' + CryptoJS.MD5(email); // chuoi ma hoa MD5
            //ma hoa mot chieu giup giu email va cho phep su dung nhu mot dinh danh
        }
    }
});
