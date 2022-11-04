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
</script>

<main><div id="navigation">
    {#each matches as match}
      <MatchButton
        {match}
        selected={match.filename === selected}
        on:click={() => selectMatch(match.filename)}
        on:keypress={() => selectMatch(match.filename)}
      />
    {/each}
  </div><span style="opacity:0;">Why do i need this?</span> <div id="main">
    {#if matchinfo.Players}
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
          ><td class="center">{matchinfo.map}</td><td class="center">{matchinfo.Duration} minutes</td><td
            class="center"
            >{selected.split("-").slice(1).join("-").slice(0, 10)}</td
          ></tr
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
          ><td
            class="tableside"
            rowspan={sortPlayers(matchinfo.Players, "CT").length + 1}
            >CT<br />{matchinfo.CT_Score}</td
          ></tr
        >
        {#each sortPlayers(matchinfo.Players, "CT") as player}
          <tr class="ct">
            <td>{player.Name}</td><td class="center">{player.Kills}</td><td
              class="center">{player.Assists}</td
            ><td class="center">{player.Deaths}</td>
            <td class="center">{player.Score}</td>
          </tr>
        {/each}
        <tr class="terrorist">
          <td
            class="tableside"
            rowspan={sortPlayers(matchinfo.Players, "CT").length + 1}
            >T<br /> {matchinfo.T_Score}</td
          ><td><br /></td>
        </tr>
        {#each sortPlayers(matchinfo.Players, "TERRORIST") as player}
          <tr class="terrorist">
            <td>{player.Name}</td><td class="center">{player.Kills}</td><td
              class="center">{player.Assists}</td
            ><td class="center">{player.Deaths}</td>
            <td class="center">{player.Score}</td>
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
