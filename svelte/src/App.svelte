<script>
  import MatchButton from "./MatchButton.svelte";
  let matches = [];
  let matchinfo = {};
  let selected = "";

  fetch(`/api/matches`).then((response) => {
    response.json().then((json) => {
      matches = json;
    });
  });

  function selectMatch(m) {
    selected = m;
    fetch(`/api/match/${m}/info`).then((response) => {
      response.json().then((json) => {
        matchinfo = json;
        console.log(matchinfo);
      });
    });
  }

  function sortPlayers(players, team) {
    return players
      .filter((player) => player.Team === team)
      .sort((a, b) => (a.Score < b.Score ? 1 : -1));
  }

  function getDate(name) {
    console.log(name);
    let date = name.split("t")[0];
    date = date.split("-").slice(1).join("-");
    let time = name.split("t")[1];
    time = time.split("-").slice(0, -1).join(":");
    console.log(date, time);
    return `${date} ${time}`;
  }
</script>

<main>
  <div id="navigation">
    {#each matches as match}
      <MatchButton
        {match}
        selected={match.filename === selected}
        on:click={() => selectMatch(match.filename)}
        on:keypress={() => selectMatch(match.filename)}
      />
    {/each}
  </div>
  <span style="opacity:0;">Why do i need this?</span>
  <div id="main">
    {#if matchinfo.filename}
      <a href="/api/match/{selected}/json">json</a>
      <a href="/api/match/{selected}/raw">raw</a>

      <table id="infotable">
        <tr
          ><th rowspan="2" class="logo"
            ><img src="map-icons/map_icon_{matchinfo.map}.svg" alt="" /></th
          ><th style="width: 150px;">Map</th><th style="width: 100px;"
            >Duration</th
          ><th style="width: 200px;">Date</th></tr
        >
        <tr
          ><td class="center">{matchinfo.map}</td><td class="center"
            >{matchinfo.duration} minutes</td
          ><td class="center">{getDate(selected)}</td></tr
        >
      </table>

      <table class="stats">
        <tr>
          <th /><th style="width: 20em;" /><th style="width: 3em;">K</th><th
            style="width: 3em;">A</th
          ><th style="width: 3em;">D</th>
          <th style="width: 3em;">Score</th>
        </tr>
        <tr class="ct"
          ><td class="tableside" rowspan={matchinfo.players_ct.length + 1}
            >CT<br />{matchinfo.score_ct}</td
          ></tr
        >
        {#each matchinfo.players_ct as player}
          <tr class="ct">
            <td>{player.name}</td><td class="center">{player.kills}</td><td
              class="center">{player.assists}</td
            ><td class="center">{player.deaths}</td>
            <td class="center">{player.score}</td>
          </tr>
        {/each}
        <tr class="terrorist">
          <td class="tableside" rowspan={matchinfo.players_t.length + 1}
            >T<br /> {matchinfo.score_t}</td
          ><td><br /></td>
        </tr>
        {#each matchinfo.players_t as player}
          <tr class="terrorist">
            <td>{player.name}</td><td class="center">{player.kills}</td><td
              class="center">{player.assists}</td
            ><td class="center">{player.deaths}</td>
            <td class="center">{player.score}</td>
          </tr>
        {/each}
      </table>
    {/if}
    <pre />
  </div>
</main>

<style>
  #navigation {
    width: 250px;
    position: absolute;
    top: 0;
    bottom: 0;
    background-color: #ffffff44;
  }

  #navigation {
    display: block;
  }

  #main {
    margin-left: 300px;
    color: #ccc;
  }

  table.stats {
    border-collapse: collapse;
  }
  .center {
    text-align: center;
  }

  .tableside {
    width: 5em;
    text-align: center;
    font-size: 150%;
  }

  .ct {
    color: #b5d4ee;
  }

  .terrorist {
    color: #ead18a;
  }

  tr.terrorist,
  tr.ct {
    border-bottom: 2px solid #ffffff11;
  }

  #infotable {
    margin-left: 100px;
    margin-bottom: 50px;
  }

  img {
    width: 50px;
  }

  a {
    color: #ffffffcc;
    font-size: 80%;
    text-decoration: none;
    border: 1px solid #ffffff55;
    border-radius: 3px;
    padding: 2px;
  }
</style>
