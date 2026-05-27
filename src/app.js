(async () => {
  const html = document.documentElement;
  const savedTheme = localStorage.getItem('theme');
  const currentTheme = savedTheme || (window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
  html.setAttribute('data-theme', currentTheme);
  document.getElementById('theme-icon').className = currentTheme === 'dark' ? 'fas fa-sun' : 'fas fa-moon';

  document.getElementById('theme-toggle').addEventListener('click', () => {
    const next = html.getAttribute('data-theme') === 'dark' ? 'light' : 'dark';
    html.setAttribute('data-theme', next);
    localStorage.setItem('theme', next);
    document.getElementById('theme-icon').className = next === 'dark' ? 'fas fa-sun' : 'fas fa-moon';
  });

  const VERSION = '2.0';
  document.getElementById('version').innerHTML = 'v' + VERSION;

  const APP_TEMPLATE = document.getElementById('app-item-template').innerHTML;
  const RESULTS = document.getElementById('results');
  const DONATIONS = document.getElementById('donations');
  const TERM = document.getElementById('term');
  const APPS = new Map();

  const modalEl = document.getElementById('json-modal');
  document.getElementById('json-modal-close').addEventListener('click', () => modalEl.classList.remove('is-active'));
  modalEl.querySelector('.modal-background').addEventListener('click', () => modalEl.classList.remove('is-active'));
  const copyBtn = document.getElementById('json-modal-copy');
  const copyLabel = copyBtn.querySelector('span:last-child');
  copyBtn.addEventListener('click', async () => {
    await navigator.clipboard.writeText(document.getElementById('json-modal-content').textContent);
    const original = copyLabel.textContent;
    copyLabel.textContent = 'Copied!';
    setTimeout(() => { copyLabel.textContent = original; }, 1200);
  });

  const doSearch = () => {
    RESULTS.innerHTML = 'Searching ...Please wait';
    DONATIONS.classList.add('is-hidden');

    const term = document.getElementById('term').value;
    const country = document.getElementById('country').value;
    const media = document.querySelector('input[name="media"]:checked').value;

    const display = (parsed) => {
      if (parsed.results && parsed.results.length > 0) {
        RESULTS.innerHTML = '';
        APPS.clear();
        DONATIONS.classList.remove('is-hidden');

        RESULTS.insertAdjacentHTML("beforeend", '<span class="tag is-link is-light is-medium">Found ' + parsed.results.length +
          (parsed.results.length === 1 ? ' entry' : ' entries') + '</span><br><br>');

        for (const app of parsed.results) {
          RESULTS.insertAdjacentHTML("beforeend", getAppHtml(APP_TEMPLATE, app));
          APPS.set(String(app.trackId || app.collectionId), app);
        }

        const els = document.querySelectorAll('.delete-app');
        for (const el of els) {
          el.addEventListener('click',
            (event) => document.getElementById('app-' + event.target.dataset.id).remove()
          );
        }
        const modal = document.getElementById('json-modal');
        const modalContent = document.getElementById('json-modal-content');
        for (const el of document.querySelectorAll('.json-app')) {
          el.addEventListener('click', (event) => {
            const id = event.currentTarget.dataset.id;
            modalContent.textContent = JSON.stringify(APPS.get(id), null, 2);
            modal.classList.add('is-active');
          });
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
        DONATIONS.classList.add('is-hidden');
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

    search(term, country, media, (err, json) => {
      if (err != null) {
        RESULTS.innerHTML = err;
        DONATIONS.classList.add('is-hidden');
      } else {
        display(json);
      }
    });
  };

  document.getElementById('search').onclick = () => doSearch();
  TERM.onkeyup = async (event) => event.key === 'Enter' && doSearch();
  document.getElementById('clear').onclick = () => {
    RESULTS.innerHTML = 'No results found';
    DONATIONS.classList.add('is-hidden');
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
