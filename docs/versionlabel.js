import ExecutionEnvironment from '@docusaurus/ExecutionEnvironment';

if (ExecutionEnvironment.canUseDOM) {

    function readVersionFile() {
        return fetch('https://raw.githubusercontent.com/stonith404/pocket-id/refs/heads/main/.version')
          .then(response => response.text())
          .catch(error => `Error reading version file: ${error}`);
      }

    function getVersion() {
        console.log("runnding")
        readVersionFile().then(version => {
          const navbarItem = document.querySelector('.navbar__item[href="#version"]');
          if (navbarItem) {
            navbarItem.innerHTML = version;
          }
        }).catch(error => console.error('Error fetching version:', error));
      }
    window.addEventListener("load", getVersion);
}