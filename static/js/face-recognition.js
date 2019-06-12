const video =  document.getElementById("video");
// const displaySize = {width: video.width, height: video.height};
var displaySize = {width: window.innerWidth, height: window.innerHeight};
const videoContainer = document.getElementById("video-container");


//download in parallel
Promise.all([
    //normal face detector, just tiny..
    faceapi.nets.tinyFaceDetector.loadFromUri('/static/js/models'),
    //mouth, nose, ect.
    faceapi.nets.faceLandmark68Net.loadFromUri('/static/js/models'),
    //allows recognize face in the box surrounded
    faceapi.nets.faceRecognitionNet.loadFromUri('/static/js/models'),
    //happy, sad, smailing etc...
    faceapi.nets.faceExpressionNet.loadFromUri('/static/js/models'),

]).then(startCamera)

function startCamera() {
    if(video != null) {


        navigator.getUserMedia(
            {video: {}},
            stream => video.srcObject = stream,
            err => console.log(err)
        )


        //video display size is different from video element size => recompute
        video.addEventListener( "loadedmetadata", function (e) {
            var videoRatio = video.videoWidth / video.videoHeight;
            // The width and height of the video element
            var width = video.offsetWidth, height = video.offsetHeight;
            // The ratio of the element's width to its height
            var elementRatio = width/height;
            // If the video element is short and wide
            if(elementRatio > videoRatio) width = height * videoRatio;
            // It must be tall and thin, or exactly equal to the original ratio
            else height = width / videoRatio;
            displaySize.width = width;
            displaySize.height = height;
        }, false );


        video.addEventListener('play', ()=> {
            const socket = websocks.GetSocket();
            const canvas = faceapi.createCanvasFromMedia(video);

            //append rectangle with face to video
            document.body.append(canvas);
            videoContainer.appendChild(canvas);
            //video.appendChild(canvas);
            faceapi.matchDimensions(canvas,displaySize);
            //detect face asynchronously from video in interval 300ms
            setInterval(async () => {
                //detects all faces from video
                const detections = await faceapi.detectAllFaces(video,
                    new faceapi.TinyFaceDetectorOptions()) //no params because emtpy default options working well
                    .withFaceLandmarks()
                    .withFaceExpressions();

                //gets detections with resized face
                const resizedDetections = faceapi.resizeResults(detections, displaySize);

                if (Array.isArray(resizedDetections) && resizedDetections.length) {
                    socket.send(JSON.stringify(
                        {
                            score: resizedDetections[0].detection.score,
                            mood: resizedDetections[0].expressions
                        }
                    ))
                }

                //clear all in canvas (rectangles around faces)
                canvas.getContext('2d').clearRect(0,0, canvas.width, canvas.height);

                //draws detections is human face for 85%
                faceapi.draw.drawDetections(canvas,resizedDetections);
                //draws happy, sad etc...
                faceapi.draw.drawFaceExpressions(canvas, resizedDetections);
            },300) //second argument says interal for recognition

        });
    }
}






