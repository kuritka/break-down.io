const video =  document.getElementById("video");
const displaySize = {width: video.width, height: video.height};


//download in parallel
Promise.all([
    //normal face detector, just tiny..
    faceapi.nets.tinyFaceDetector.loadFromUri('/static/js/models'),
    //mouth, nose, ect.
    faceapi.nets.faceLandmark68Net.loadFromUri('/static/js/models'),
    //allows recognize face in the box surrounded
    faceapi.nets.faceRecognitionNet.loadFromUri('/static/js/models'),
    //happy, sad, smailing etc...
    faceapi.nets.faceExpressionNet.loadFromUri('/static/js/models')
]).then(startCamera)

function startCamera() {
    navigator.getUserMedia(
        {video: {}},
        stream => video.srcObject = stream,
        err => console.log(err)
    )
}

video.addEventListener('play', ()=> {
    const canvas = faceapi.createCanvasFromMedia(video);
    //append rectangle with face to video
    console.log(video.style)
    document.body.append(canvas);
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
        //clear all in canvas (rectangles around faces)
        canvas.getContext('2d').clearRect(0,0, canvas.width, canvas.height);

        //draws detections is human face for 85%
        faceapi.draw.drawDetections(canvas,resizedDetections);
        //draws happy, sad etc...
        faceapi.draw.drawFaceExpressions(canvas, resizedDetections);
    },300) //second argument says interal for recognition

})

