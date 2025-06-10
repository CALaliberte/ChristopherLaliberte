function Landing() {
  document.querySelector(".Biography").classList.remove("show");
  document.querySelector(".Biography").classList.add("hide");
  document.querySelector(".Scores").classList.remove("show");
  document.querySelector(".Scores").classList.add("hide");
  document.querySelector(".Contact").classList.remove("show");
  document.querySelector(".Contact").classList.add("hide");
  document.querySelector(".Landing").classList.add("show");
}

function Biography() {
  document.querySelector(".Landing").classList.remove("show");
  document.querySelector(".Landing").classList.add("hide");
  document.querySelector(".Scores").classList.remove("show");
  document.querySelector(".Scores").classList.add("hide");
  document.querySelector(".Contact").classList.remove("show");
  document.querySelector(".Contact").classList.add("hide");
  document.querySelector(".Biography").classList.add("show");
}

function Scores() {
  document.querySelector(".Landing").classList.remove("show");
  document.querySelector(".Landing").classList.add("hide");
  document.querySelector(".Biography").classList.remove("show");
  document.querySelector(".Biography").classList.add("hide");
  document.querySelector(".Contact").classList.remove("show");
  document.querySelector(".Contact").classList.add("hide");
  document.querySelector(".Scores").classList.add("show");
}

function Contact() {
  document.querySelector(".Landing").classList.remove("show");
  document.querySelector(".Landing").classList.add("hide");
  document.querySelector(".Biography").classList.remove("show");
  document.querySelector(".Biography").classList.add("hide");
  document.querySelector(".Scores").classList.remove("show");
  document.querySelector(".Scores").classList.add("hide");
  document.querySelector(".Contact").classList.add("show");
}

document.addEventListener("DOMContentLoaded", function () {
  const jumpToTopButton = document.getElementById("jumpToTop");

  window.onscroll = function () {
    if (document.documentElement.scrollTop > 100 || document.body.scrollTop > 100) {
      jumpToTopButton.style.display = "block";
    } else {
      jumpToTopButton.style.display = "none";
    }
  };

  window.jumpToTop = function () {
    document.documentElement.scrollTop = 0;
    document.body.scrollTop = 0;
  };
});

function openPDF(filePath) {
  const pdfViewer = document.getElementById('pdfViewer');
  const pdfFrame = document.getElementById('pdfFrame');
  pdfViewer.showModal();
  pdfFrame.src = filePath;

  pdfViewer.addEventListener('click', (event) => {
    if (event.target === pdfViewer) {
      closePDF();
    }
  });
}

function closePDF() {
  const pdfViewer = document.getElementById('pdfViewer');
  pdfViewer.close();
}

function openLink(linkURL) {
  var popupWindow = window.open(linkURL, "_blank", "noopener");
  popupWindow.name = "pdfPopup_" + Date.now();
}

document.addEventListener('play', function(e){  
  var audios = document.getElementsByTagName('audio');  
  for(var i = 0, len = audios.length; i < len;i++){  
      if(audios[i] != e.target){  
          audios[i].pause();
          audios[i].currentTime = 0;  
    }
  }
}, true);

const API_KEY = 'AIzaSyCyv-VFWCylVvwORP_jBkeuI7e3yQE-7ho';
const CHANNEL_ID = 'UCXxyDExQviakr4tr3HB33hQ';

document.addEventListener('DOMContentLoaded', fetchLatestVideo);

function fetchLatestVideo() {
  const searchUrl = `https://www.googleapis.com/youtube/v3/search?part=snippet&channelId=${CHANNEL_ID}&order=date&type=video&maxResults=5&key=${API_KEY}`;

  fetch(searchUrl)
    .then((response) => response.json())
    .then((searchData) => {
      if (searchData.items && searchData.items.length > 0) {
        const videoIds = searchData.items.map((item) => item.id.videoId).join(',');
        
        const videoUrl = `https://www.googleapis.com/youtube/v3/videos?part=snippet,status&id=${videoIds}&key=${API_KEY}`;
        return fetch(videoUrl);
      } else {
        throw new Error('No videos found for this channel.');
      }
    })
    .then((response) => response.json())
    .then((videoData) => {
      if (videoData.items && videoData.items.length > 0) {
        const sortedVideos = videoData.items.sort((a, b) => {
          const dateA = new Date(a.snippet.publishedAt);
          const dateB = new Date(b.snippet.publishedAt);
          return dateB - dateA;
        });

        const latestVideo = sortedVideos[0];
        const videoId = latestVideo.id;
        const videoTitle = latestVideo.snippet.title;
        const publishedAt = latestVideo.snippet.publishedAt;

        console.log(`Latest Public Video: ${videoTitle} (Published at: ${publishedAt})`);

        const iframe = document.getElementById('youtube-video');
        iframe.src = `https://www.youtube.com/embed/${videoId}`;
      } else {
        throw new Error('No public videos found.');
      }
    })
    .catch((error) => console.error('Error fetching video:', error));
}
