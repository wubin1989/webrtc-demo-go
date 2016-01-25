'use strict';

// navigator.webkitGetUserMedia({video: true, audio: true}, function(stream){
	
	// var ws = io.connect('http://192.168.42.44:4000');

	// ws.on('connect', function(){
	// 	console.log(ws);
	// });

	// ws.on('returned signal', function(data){
	// 	document.getElementById('otherId').value = data;
	// });

	var room = location.search && location.search.split('?')[1];

    console.log(room);

	var webrtc = new SimpleWebRTC({
        // the id/element dom element that will hold "our" video
        localVideoEl: 'localVideo',
        // the id/element dom element that will hold remote videos
        remoteVideosEl: '',
        // immediately ask for camera access
		autoRequestMedia: true,
		debug: false,
		detectSpeakingEvents: true,
		autoAdjustMic: false,
        url: 'http://localhost:8451'
    });

    webrtc.on('readyToCall', function (socketId) {
        console.log('socketId is %s', socketId);
        // you can name it anything
        if (room) {
            console.log('room is %s', room);
            webrtc.joinRoom(room);
        }
    });

    // function showVolume(el, volume) {
    //     if (!el) return;
    //     if (volume < -45) { // vary between -45 and -20
    //         el.style.height = '0px';
    //     } else if (volume > -20) {
    //         el.style.height = '100%';
    //     } else {
    //         el.style.height = '' + Math.floor((volume + 100) * 100 / 25 - 220) + '%';
    //     }
    // }
    // webrtc.on('channelMessage', function (peer, label, data) {
    //     if (data.type == 'volume') {
    //         showVolume(document.getElementById('volume_' + peer.id), data.volume);
    //     }
    // });

    webrtc.on('videoAdded', function (video, peer) {
        console.log('video added', peer);
        var remotes = document.getElementById('remotes');
        if (remotes) {
            debugger;
            var d = document.createElement('div');
            d.className = 'videoContainer';
            d.id = 'container_' + webrtc.getDomId(peer);
            d.appendChild(video);
            var vol = document.createElement('div');
            vol.id = 'volume_' + peer.id;
            vol.className = 'volume_bar';
            video.onclick = function () {
                video.style.width = video.videoWidth + 'px';
                video.style.height = video.videoHeight + 'px';
            };
            d.appendChild(vol);
            remotes.appendChild(d);
        }
    });

    webrtc.on('videoRemoved', function (video, peer) {
        console.log('video removed ', peer);
        var remotes = document.getElementById('remotes');
        var el = document.getElementById('container_' + webrtc.getDomId(peer));
        if (remotes && el) {
            remotes.removeChild(el);
        }
    });

    // webrtc.on('volumeChange', function (volume, treshold) {
    //     //console.log('own volume', volume);
    //     showVolume(document.getElementById('localVolume'), volume);
    // });

    // Since we use this twice we put it here
    // function setRoom(name) {
    //     $('form').remove();
    //     $('h1').text(name);
    //     $('#subTitle').text('Link to join: ' + location.href);
    //     $('body').addClass('active');
    // }

    if (room) {
        // setRoom(room);
    } else {
        console.log(123);
        //$('form').submit(function () {
            //var val = $('#sessionInput').val().toLowerCase().replace(/\s/g, '-').replace(/[^A-Za-z0-9_\-]/g, '');
            var val = 'chat';
            webrtc.createRoom(val, function (err, name) {
                console.log(' create room cb', arguments);
            });

            //return false;          
        //});
    }

    // var button = $('#screenShareButton'),
    //     setButton = function (bool) {
    //         button.text(bool ? 'share screen' : 'stop sharing');
    //     };
    // webrtc.on('localScreenStopped', function () {
    //     setButton(true);
    // });

    // setButton(true);

    // button.click(function () {
    //     if (webrtc.getLocalScreen()) {
    //         webrtc.stopScreenShare();
    //         setButton(true);
    //     } else {
    //         webrtc.shareScreen(function (err) {
    //             if (err) {
    //                 setButton(true);
    //             } else {
    //                 setButton(false);
    //             }
    //         });
            
    //     }
    // });


// }, function(err){
// 	console.log(err);
// })


	

