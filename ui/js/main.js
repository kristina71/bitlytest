document.addEventListener('DOMContentLoaded', function () {
  document.getElementById("preloader").classList.remove("active");
  var xhr = new XMLHttpRequest();

  xhr.open('POST', '/all', false);
  xhr.send();
  if (xhr.status != 200) {
    alert(xhr.status + xhr.responseText);
  } else {
    var obj = JSON.parse(xhr.responseText);

    for (i = 0; i < obj.length; i++) {
      document.getElementById("app").innerHTML += "<div class=\"row\">" +
        "<form method=\"POST\" action=\"/edit\" id=\"form\">"+
        "<input type=\"hidden\" name=\"id\" value=\""+obj[i].id+"\">"+
        "<div class=\"input-field col s2\"><label>http://"+ document.location.host+"/</label></div>"+
        "<div class=\"input-field col s3\">"+
        "<input type=\"text\" name=\"small_url\" value=\"" + obj[i].small_url + "\">" +
        "</div><div class=\"input-field col s0.5\">=&gt;</div><div class=\"input-field col s3\"><input type=\"text\" name=\"origin_url\" value=\"" + obj[i].origin_url + "\"></div>" +
        "<div class=\"input-field col s1\"> <input type=\"submit\" name=\"save\" value=\"Save\" class=\"waves-effect waves-light btn\">" +
        "</div></form>"+
        "<form method=\"POST\" action=\"/delete\" id=\"form\">"+
        "<div class=\"input-field col s1\">"+
        "<input type=\"hidden\" name=\"id\" value=\""+obj[i].id+"\">"+
        "<input type=\"submit\" class=\"waves-effect waves-light btn\" value=\"X\"></div></form>" +
        "<a href=\"http://"+ document.location.host+"/"+obj[i].small_url+"\">Open</a></div>";
    }
  }
});
document.getElementById("app").onerror = function () {
  alert("Something went wrong");
};