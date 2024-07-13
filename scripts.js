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
