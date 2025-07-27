<script>
  import { onMount } from 'svelte';
  import 'leaflet/dist/leaflet.css'; // Import Leaflet's CSS
  import * as L from 'leaflet'; // Import Leaflet's JavaScript library

  let stations = [];
  let error = null;
  let loading = true;

  // REMOVE THE customMarkerIcon DEFINITION - it's no longer needed for circle markers.
  // const customMarkerIcon = L.icon({
  //   iconUrl: '/images/reddot.png',
  //   shadowUrl: '/images/reddot.png',
  //   iconSize:    [25, 41],
  //   iconAnchor:  [12, 41],
  //   popupAnchor: [1, -34],
  //   shadowSize:  [41, 41]
  // });

  onMount(async () => {
    try {
      const response = await fetch('http://localhost:8080/api/stations');
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      stations = await response.json();

      setTimeout(() => {
        stations.forEach(station => {
          const mapDivId = `map-${station.StationID}`;
          if (document.getElementById(mapDivId)) {
            const map = L.map(mapDivId).setView([station.CenterLatitude, station.CenterLongitude], 10);

            L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
              attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
            }).addTo(map);

            // ✨ CHANGE THIS LINE: Use L.circleMarker instead of L.marker with an icon
            L.circleMarker([station.CenterLatitude, station.CenterLongitude], {
              radius: 6,         // Size of the circle
              fillColor: 'red',  // Fill color of the circle
              color: 'darkred',  // Border color of the circle
              weight: 1,         // Border width
              opacity: 1,        // Opacity of the border
              fillOpacity: 0.8   // Opacity of the fill
            }).addTo(map);

          }
        });
      }, 0);
    } catch (e) {
      error = e.message;
      console.error("Failed to fetch stations:", e);
    } finally {
      loading = false;
    }
  });
</script>

<main>
  <h1>🌊 Tidal Station Data</h1>

  {#if loading}
    <p>Loading data...</p>
  {:else if error}
    <p style="color: red;">Error: {error}</p>
  {:else if stations.length > 0}
    <div class="station-container">
      {#each stations as station (station.StationID)}
        <div class="station-card">
          <div>
            <h2>{station.StationName} ({station.StationNameLocal})</h2>
            <p><strong>ID:</strong> {station.Id}</p>
            <p><strong>Station ID:</strong> {station.StationID}</p>
            <p><strong>Coordinates:</strong> {station.CenterLatitude}, {station.CenterLongitude}</p>
            <p><strong>Dates:</strong> {station.EarliestDate} to {station.LatestDate}</p>
            <p><a href={station.AccessURL} target="_blank" rel="noopener noreferrer">Access Data</a></p>
          </div>
          <div id="map-{station.StationID}" class="map-container"></div>
        </div>
      {/each}
    </div>
  {:else}
    <p>No station data available.</p>
  {/if}
</main>

<style>
  @import 'leaflet/dist/leaflet.css';

  main {
    font-family: sans-serif;
    padding: 20px;
    max-width: 1200px;
    margin: 0 auto;
    background-color: #f4f4f4;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0,0,0,0.1);
  }

  h1 {
    text-align: center;
    color: #333;
    margin-bottom: 30px;
  }

  .station-container {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(550px, 1fr));
    gap: 20px;
  }

  .station-card {
    background-color: #fff;
    border: 1px solid #ddd;
    border-radius: 8px;
    padding: 20px;
    box-shadow: 0 1px 3px rgba(0,0,0,0.08);
    transition: transform 0.2s ease-in-out;
    display: flex;
    align-items: flex-start;
    gap: 20px;
  }

  .station-card:hover {
    transform: translateY(-5px);
  }

  .station-card > div:first-child {
    flex: 1;
  }

  .station-card h2 {
    color: #0056b3;
    margin-top: 0;
    font-size: 1.3em;
    margin-bottom: 10px;
  }

  .station-card p {
    margin: 5px 0;
    color: #555;
    line-height: 1.5;
  }

  .station-card a {
    color: #007bff;
    text-decoration: none;
  }

  .station-card a:hover {
    text-decoration: underline;
  }

  .map-container {
    width: 200px;
    height: 150px;
    border: 1px solid #ccc;
    border-radius: 4px;
    flex-shrink: 0;
  }
</style>