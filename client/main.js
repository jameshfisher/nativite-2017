var audio = new Audio('./lamb.mp3');

var nameToVec = {
  sophie:    { x:  0, y: -0.7 },
  constance: { x: -1, y:  0   },
  victoire:  { x:  1, y:  0   },
  felicite:  { x:  0, y:  0.7 }
};

var nameToLabelOffset = {
  sophie:    { x: 0,  y: 0 },
  constance: { x: -2, y: 0 },
  victoire:  { x: 3,  y: 0 },
  felicite:  { x: 0,  y: 0 }
};

var nameToScoreOffset = {
  sophie:    { x:  -14, y: 6 },
  constance: { x: -16, y: 6 },
  victoire:  { x:  17, y: 6 },
  felicite:  { x:  -14, y: 9 }
};

function showScore(fille, score) {
  var distSteps = 3 + score;
  var distPx = distSteps * 25;
  console.log("Setting score for", fille, "to", score);
  var vec = nameToVec[fille];
  var leftPx = Math.round(vec.x*distPx);
  var topPx = Math.round(vec.y*distPx);
  var sheep = document.getElementById(fille);
  if (sheep) {
    sheep.style.marginLeft = (leftPx-32) + 'px';
    sheep.style.marginTop  = (topPx -32) + 'px';
  } else {
    console.log("Could not find sheep for", fille);
  }
  var label = document.getElementById("label_" + fille);
  if (label) {
    label.style.marginLeft = ((leftPx+nameToLabelOffset[fille].x*4)-72) + "px";
    label.style.marginTop  = ((topPx+nameToLabelOffset[fille].y*4) -96) + "px";
  } else {
    console.log("Could not find label for", fille);
  }
  var score_el = document.getElementById("score_" + fille);
  if (score_el) {
    score_el.style.backgroundPosition = "0px" + " " + (-24*score) + "px";
    score_el.style.marginLeft = ((leftPx+nameToScoreOffset[fille].x*4)-32) + "px";
    score_el.style.marginTop  = ((topPx+nameToScoreOffset[fille].y*4) -32) + "px";
  } else {
    console.log("Could not find score element for", fille);
  }
}

function showScores(snapshot) {
  let winner = false;
  for (fille in snapshot) {
    const score = snapshot[fille];
    if (score == 0) winner = true;
    showScore(fille, score);
  }

  document.getElementById("snow").style.opacity = winner ? "1": "0";
}

fetch("https://nativite-2017.herokuapp.com/scores").then(function(response) {
  return response.json();
}).then(showScores);

// FIXME do this from Pusher
function onUpdateSnapshot(snapshot) {
  console.log("Got value", snapshot);
  audio.currentTime = 0;
  audio.play();
  showScores(snapshot);
}
