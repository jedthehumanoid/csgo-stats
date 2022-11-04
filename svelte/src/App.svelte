<script>
  let matches = [];
  let matchinfo = {};
  let ct = 0;
  let t = 0;
  fetch(`/api/matches`).then((response) => {
    response.json().then((json) => {
      matches = json;
    });
  });

  function selectMatch(m) {
    fetch(`/api/matchinfo/${m}`).then((response) => {
      response.json().then((json) => {
        matchinfo = json;
        console.log(matchinfo);
      });
    });
  }
</script>

<main>
  <div id="navigation">
    {#each matches as match}
      <div class="match" on:click={() => selectMatch(match.filename)}>
        <img
          class="floatleft"
          src="map-icons/map_icon_{match.map}.svg"
          style="width: 50px;"
        />
        <div class="floatleft">
          <span class="ct">{match.ct_score}</span> - <span class="terrorist">{match.t_score}</span><br
          />{match.filename.slice(0, 10)}
        </div>
      </div>
    {/each}
  </div>
  <div id="main">
    {#if matchinfo.Players}
      Map: {matchinfo.map}
      <table class="ct">
        <tr>
          <th
            class="tableside"
            rowspan={matchinfo.Players.filter((player) => player.Team === "CT")
              .length + 1}>CT<br />{matchinfo.CT_Score}</th
          ><th style="width: 20em;" /><th style="width: 4em;">K</th><th
            style="width: 4em;">A</th
          ><th style="width: 4em;">D</th>
          <th style="width: 4em;">Score</th>
        </tr>
        {#each matchinfo.Players.filter((player) => player.Team === "CT").sort( (a, b) => (a.Score < b.Score ? 1 : -1) ) as player}
          <tr>
            <td>{player.Name}</td><td class="center">{player.Kills}</td><td
              class="center">{player.Assists}</td
            ><td class="center">{player.Deaths}</td>
            <td class="center">{player.Score}</td>
          </tr>
        {/each}
      </table>
      <div class="spacer" />
      <table class="terrorist">
        <tr>
          <th
            class="tableside"
            rowspan={matchinfo.Players.filter(
              (player) => player.Team === "TERRORIST"
            ).length + 1}>T<br /> {matchinfo.T_Score}</th
          ><th style="width: 20em;" /><th style="width: 4em;" /><th
            style="width: 4em;"
          /><th style="width: 4em;" />
          <th style="width: 4em;" />
        </tr>
        {#each matchinfo.Players.filter((player) => player.Team === "TERRORIST").sort( (a, b) => (a.Score < b.Score ? 1 : -1) ) as player}
          <tr>
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
    padding: 10px;
  }

  #navigation div {
    display: block;
  }

  #main {
    margin-left: 300px;
    color: #ccc;
    padding: 75px;
  }

  .match {
    cursor: pointer;
    height: 75px;
    border-bottom: 1px solid #000000;
  }

  table {
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

  .spacer {
    height: 1em;
  }
  .ct {
    color: #b5d4ee;
  }

  .terrorist {
    color: #ead18a;
  }

  .floatleft {
    float: left;
  }

  tr {
    border-bottom: 2px solid #ffffff11;
  }
</style>
