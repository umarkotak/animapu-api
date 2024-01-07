OTAKUDESU SCRAP!

Phase 1
1. Open target anime episode: https://otakudesu.wiki/episode/mt-ithd-s2-episode-1-sub-indo/
2. Find <link rel='shortlink' href='https://otakudesu.wiki/?p=146127' />
3. Find window.__x__nonce ? $.ajax("https://otakudesu.wiki/wp-admin/admin-ajax.php", { with action 2a3505c93b0035d3f455df82bf976b84
4. <iframe src="https://desustream.me/beta/stream/?id=SVRBZldYQWFvSVBjUDd3eEhlR21Tc0VqdzBxZmxUTGFVZlF6Ulh6QXZVTT0=" WIDTH="420" HEIGHT="370" allowfullscreen="true" webkitallowfullscreen="true" mozallowfullscreen="true"></iframe>
5. Get the iframe value SVRBZldYQWFvSVBjUDd3eEhlR21Tc0VqdzBxZmxUTGFVZlF6Ulh6QXZVTT0

Phase 2
1. API call to get the nonce
curl 'https://otakudesu.wiki/wp-admin/admin-ajax.php' \
  -H 'authority: otakudesu.wiki' \
  -H 'accept: */*' \
  -H 'accept-language: en-US,en;q=0.9,id;q=0.8' \
  -H 'content-type: application/x-www-form-urlencoded; charset=UTF-8' \
  -H 'cookie: _ga=GA1.2.1861696737.1695554920; _gid=GA1.2.526398125.1696081168; _gat=1; _ga_025LZFQCB2=GS1.2.1696081167.2.1.1696082456.0.0.0' \
  -H 'origin: https://otakudesu.wiki' \
  -H 'referer: https://otakudesu.wiki/episode/mt-ithd-s2-episode-1-sub-indo/' \
  -H 'sec-ch-ua: "Google Chrome";v="117", "Not;A=Brand";v="8", "Chromium";v="117"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "Windows"' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-site: same-origin' \
  -H 'user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36' \
  -H 'x-requested-with: XMLHttpRequest' \
  --data-raw 'action=aa1208d27f29ca340c92c66d1926f13f' \
  --compressed
  {"data":"95d72c2124"}
2. API call to get data
curl 'https://otakudesu.wiki/wp-admin/admin-ajax.php' \
  -H 'authority: otakudesu.wiki' \
  -H 'accept: */*' \
  -H 'accept-language: en-US,en;q=0.9,id;q=0.8' \
  -H 'content-type: application/x-www-form-urlencoded; charset=UTF-8' \
  -H 'cookie: _ga=GA1.2.1861696737.1695554920; _gid=GA1.2.526398125.1696081168; _gat=1; _ga_025LZFQCB2=GS1.2.1696081167.2.1.1696082456.0.0.0' \
  -H 'origin: https://otakudesu.wiki' \
  -H 'referer: https://otakudesu.wiki/episode/mt-ithd-s2-episode-1-sub-indo/' \
  -H 'sec-ch-ua: "Google Chrome";v="117", "Not;A=Brand";v="8", "Chromium";v="117"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "Windows"' \
  -H 'sec-fetch-dest: empty' \
  -H 'sec-fetch-mode: cors' \
  -H 'sec-fetch-site: same-origin' \
  -H 'user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36' \
  -H 'x-requested-with: XMLHttpRequest' \
  --data-raw 'id=146127&i=0&q=720p&nonce=95d72c2124&action=2a3505c93b0035d3f455df82bf976b84' \
  --compressed
  {
    "data": "PGRpdiBjbGFzcz0icmVzcG9uc2l2ZS1lbWJlZC1zdHJlYW0iPjxpZnJhbWUgc3JjPSJodHRwczovL2Rlc3VzdHJlYW0ubWUvYmV0YS9zdHJlYW0vaGQvP2lkPVNWUkJabGRZUVdGdlNWQmpVRGQzZUVobFIyMVRjMFZxZHpCeFpteFVUR0ZWWmxGNlVsaDZRWFpWVFQwPSIgV0lEVEg9IjQyMCIgSEVJR0hUPSIzNzAiIGFsbG93ZnVsbHNjcmVlbj0idHJ1ZSIgd2Via2l0YWxsb3dmdWxsc2NyZWVuPSJ0cnVlIiBtb3phbGxvd2Z1bGxzY3JlZW49InRydWUiPjwvaWZyYW1lPjwvZGl2Pg=="
  }
  use base64 decode to decode the data and get the iframe element src
3. https://desustream.me/beta/stream/ + hd + /?id=SVRBZldYQWFvSVBjUDd3eEhlR21Tc0VqdzBxZmxUTGFVZlF6Ulh6QXZVTT0= this will result on mp4 obj
curl 'https://desustream.me/beta/stream/hd/?id=SVRBZldYQWFvSVBjUDd3eEhlR21Tc0VqdzBxZmxUTGFVZlF6Ulh6QXZVTT0=' \
  -H 'authority: desustream.me' \
  -H 'accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7' \
  -H 'accept-language: en-US,en;q=0.9,id;q=0.8' \
  -H 'sec-ch-ua: "Google Chrome";v="117", "Not;A=Brand";v="8", "Chromium";v="117"' \
  -H 'sec-ch-ua-mobile: ?0' \
  -H 'sec-ch-ua-platform: "Windows"' \
  -H 'sec-fetch-dest: iframe' \
  -H 'sec-fetch-mode: navigate' \
  -H 'sec-fetch-site: cross-site' \
  -H 'sec-fetch-user: ?1' \
  -H 'upgrade-insecure-requests: 1' \
  -H 'user-agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36' \
  --compressed
  sources: [
  {'file':'https://rr2---sn-npoe7ne6.googlevideo.com/videoplayback?expire=1696111269&ei=JSoYZbbhHK-yn88PqZKogAo&ip=2a04:3543:1000:2310:30da:13ff:fead:6be6&id=3ed345110073fc1d&itag=22&source=blogger&xpc=Egho7Zf3LnoBAQ%3D%3D&mh=Lu&mm=31&mn=sn-npoe7ne6&ms=au&mv=m&mvi=2&pl=32&susc=bl&eaua=3mympojr-Ug&mime=video/mp4&vprv=1&dur=1420.062&lmt=1688931149346306&mt=1696082014&txp=1311224&sparams=expire,ei,ip,id,itag,source,xpc,susc,eaua,mime,vprv,dur,lmt&sig=AOq0QJ8wRQIhALHgxYbeem2nEll0YCk7KEcQB_G9yFhr60rvZyHZTF0wAiAXZrxs0Ag1ujV3AW5GNoweqmDaDQYlTpZf5DZCMsCyrw%3D%3D&lsparams=mh,mm,mn,ms,mv,mvi,pl&lsig=AG3C_xAwRQIgHmUOwdoUq0fWhV55ZodXvVQOCx4J5M0rugG974LX_3cCIQDwcBJ1eAuAzVzNbDBPzYJ3im9kw525qOpU1bL_iCXOwQ%3D%3D',
  'type':'video/mp4'}],
  image: "https://1.bp.blogspot.com/-_hry8Jf4-rE/YUFwSer_GmI/AAAAAAAAJzQ/x3FebfvoOmcpwVGyUmkDqx5a7t8-GmN8ACLcBGAsYHQ/s960/Otakudesu%2BStreaming.png",
  captions:
      {
      color:'#FFFF00',fontSize:17,backgroundOpacity:50
  },

<script>
    window.__x__nonce = null,
    $('.mirrorstream a[href^="#"]').on("click", function(a) {
        a.preventDefault();
        const n = a.currentTarget
          , e = JSON.parse(atob(n.dataset.content));
        window.__x__nonce ? $.ajax("https://otakudesu.wiki/wp-admin/admin-ajax.php", {
            method: "POST",
            processData: !0,
            cache: !0,
            data: {
                ...e,
                nonce: window.__x__nonce,
                action: "2a3505c93b0035d3f455df82bf976b84"
                2a3505c93b0035d3f455df82bf976b84
         }
        }).done(({data: a})=>{
            document.getElementById("pembed").innerHTML = atob(a)
        }
        ).fail(function() {}) : $.ajax("https://otakudesu.wiki/wp-admin/admin-ajax.php", {
            method: "POST",
            processData: !0,
            cache: !0,
            data: {
                action: "aa1208d27f29ca340c92c66d1926f13f"
            }
        }).done(({data: a})=>{
            window.__x__nonce = a,
            $.ajax("https://otakudesu.wiki/wp-admin/admin-ajax.php", {
                method: "POST",
                processData: !0,
                cache: !0,
                data: {
                    ...e,
                    nonce: a,
                    action: "2a3505c93b0035d3f455df82bf976b84"
                }
            }).done(({data: a})=>{
                document.getElementById("pembed").innerHTML = atob(a)
            }
            ).fail(function() {})
        }
        ).fail(function() {})
    });
</script>
