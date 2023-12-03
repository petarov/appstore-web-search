(async () => {
  document.querySelector('.controls').style.display = 'none';

  const go = new Go();
  if (WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming(fetch('appstore.wasm'), go.importObject)
      .then(result => {
        document.querySelector('.loading').style.display = 'none';
        document.querySelector('.controls').style.display = 'block';
        go.run(result.instance);
        get_app_version((v) => document.getElementById('version').innerHTML = 'v' + v);
      });
  } else {
    // Safari
    fetch('appstore.wasm').then(response => response.arrayBuffer())
      .then(bytes => WebAssembly.instantiate(bytes, go.importObject))
      .then(result => {
        document.querySelector('.loading').style.display = 'none';
        document.querySelector('.controls').style.display = 'block';
        go.run(result.instance);
        get_app_version((v) => document.getElementById('version').innerHTML = 'v' + v);
      });
  }

  const APP_TEMPLATE = document.getElementById('app-item-template').innerHTML;
  const RESULTS = document.getElementById('results');
  const TERM = document.getElementById('term');

  const doSearch = () => {
    RESULTS.innerHTML = 'Searching ...Please wait';

    const term = document.getElementById('term').value;
    const country = document.getElementById('country').value;
    const media = document.querySelector('input[name="media"]:checked').value;

    const display = (parsed) => {
      if (parsed.results && parsed.results.length > 0) {
        RESULTS.innerHTML = '';
        for (const app of parsed.results) {
          RESULTS.insertAdjacentHTML("beforeend", getAppHtml(APP_TEMPLATE, app));
        }

        const els = document.querySelectorAll('.delete-app');
        for (const el of els) {
          el.addEventListener('click',
            (event) => document.getElementById('app-' + event.target.dataset.id).remove()
          );
        }
        if (navigator.share) {
          const els = document.querySelectorAll('.share-app');
          for (const el of els) {
            el.addEventListener('click', (event) => {
              const trackId = event.target.dataset.id;
              const title = event.target.dataset.title;
              const link = event.target.dataset.link;
              (async () => {
                const data = {
                  title: title,
                  text: '[App]: ' + title + ' [iTunes ID]: ' + trackId,
                  url: link
                };
                await navigator.share(data);
              })();
            });
          }
        }
      } else {
        RESULTS.innerHTML = 'No results found';
      }
    };

    const search = (term, country, media, cb) => {
      const url = isNaN(term) ?
        `https://itunes.apple.com/search?media=${media}&term=${term}&country=${country}&callback=_cb` :
        `https://itunes.apple.com/lookup?id=${term}&country=${country}&callback=_cb`;

      fetch(url).then(response => {
        if (response.status / 100 === 2) {
          response.text().then(body => {
            body = body.substring(7, body.length - 4);
            cb(null, JSON.parse(body));
          });
        } else {
          cb(`HTTP Err: ${response.status} ${response.statusText}`, null);
        }
      }).catch(err => cb(err, null));
    };

    get_cache(term, country, media, (err, entry) => {
      if (err) {
        search(term, country, media, (err, json) => {
          if (err != null) {
            RESULTS.innerHTML = err;
            if (json != null) {
              RESULTS.innerHTML += "<br>";
              RESULTS.innerHTML += JSON.parse(json).errorMessage;
            }
          } else {
            display(json);
            put_cache(term, country, media, JSON.stringify(json), (err, key) => {
              if (err) {
                console.error('cache put failed!', err);
              }
            });
          }
        });
      } else {
        console.log('** cache hit', entry.key);
        display(JSON.parse(entry.data));
      }
    });
  };

  document.getElementById('search').onclick = () => doSearch();
  TERM.onkeyup = async (event) => event.key === 'Enter' && doSearch();
  document.getElementById('clear').onclick = () => {
    RESULTS.innerHTML = 'No results found';
    TERM.value = '';
    TERM.focus();
  };

  let country = navigator.language.substr(0, 2).toLowerCase();
  country = country == 'en' ? 'US' : country;
  document.getElementById('country').value = country;
})();

function getAppHtml(template, app) {
  var tpl = template.slice(0);
  tpl = tpl.replace(/__LINK__/g, app.trackViewUrl || app.collectionViewUrl);
  tpl = tpl.replace(/__IMG_100__/g, app.artworkUrl100);
  tpl = tpl.replace(/__IMG_60__/g, app.artworkUrl60);
  tpl = tpl.replace(/__TRACK_ID__/g, app.trackId || app.collectionId);
  if (app.bundleId) {
    tpl = tpl.replace(/__BUNDLE_ID__/g, app.bundleId);
  } else {
    tpl = tpl.replace(/__BUNDLE_ID__HIDDEN__/g, 'is-hidden');
  }
  tpl = tpl.replace(/__TITLE__/g, app.trackName || app.collectionName);
  tpl = tpl.replace(/__TYPE__/g, app.kind || app.wrapperType);
  tpl = tpl.replace(/__ARTIST__/g, app.artistName);
  if (app.fileSizeBytes) {
    tpl = tpl.replace(/__SIZE__/g, parseFloat(app.fileSizeBytes / 1024 / 1024).toFixed(2));
  } else {
    tpl = tpl.replace(/__SIZE_HIDDEN__/g, 'is-hidden');
  }
  if (app.genres) {
    tpl = tpl.replace(/__GENRES__/g, app.genres.join(', '));
  } else {
    tpl = tpl.replace(/__GENRES__/g, app.primaryGenreName);
  }
  if (app.price && app.price > 0.0) {
    tpl = tpl.replace(/__PRICE_STYLE__/g, 'is-danger');
    tpl = tpl.replace(/__PRICE__/g, app.formattedPrice);
  } else if (app.trackPrice && app.trackPrice > 0.0) {
    tpl = tpl.replace(/__PRICE_STYLE__/g, 'is-danger');
    tpl = tpl.replace(/__PRICE__/g, app.trackPrice + ' ' + app.currency);
  } else {
    tpl = tpl.replace(/__PRICE_STYLE__/g, 'is-success');
    tpl = tpl.replace(/__PRICE__/g, app.formattedPrice || 'free');
  }
  tpl = tpl.replace(/__VERSION__/g, app.version);
  if (navigator.share) {
    tpl = tpl.replace(/__SHARE_HIDDEN__/g, '');
  } else {
    tpl = tpl.replace(/__SHARE_HIDDEN__/g, 'is-hidden');
  }
  return tpl;
}
