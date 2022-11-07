<script lang="ts">
  import MatchButton from "./MatchButton.svelte";
  import type { MatchInfo, MatchBrief } from "./csgo";

  let matches: MatchBrief[] = [];

  let matchinfo: MatchInfo = {
    filename: "",
    map: "",
    time: "",
    duration: 0,
    players_ct: [],
    players_t: [],
    score_ct: 0,
    score_t: 0,
    rounds: [],
  };

  let selected: string = "";
  let page: number = 0;

  fetch(`/api/matches`).then((response) => {
    response.json().then((json) => {
      matches = json;
    });
  });

  function selectMatch(m: string): void {
    selected = m;
    fetch(`/api/match/${m}/info`).then((response) => {
      response.json().then((json) => {
        matchinfo = json;
        console.log(matchinfo);
      });
    });
  }

  function getRoundLogo(notice: string): string {
    console.log(notice);
    if (notice == "defused") {
      return "&#x2702;&#xFE0E;";
    }
    if (notice == "bombed") {
      return "&#x1F4A5;&#xFE0E;";
    }
    return "&#x1F480;&#xFE0E;";
  }
</script>

<main>
  <div id="navigation">
    {#each matches.slice(page * 15, (page + 1) * 15) as match}
      <MatchButton
        {match}
        selected={match.filename === selected}
        on:click={() => selectMatch(match.filename)}
        on:keypress={() => selectMatch(match.filename)}
      />
    {/each}
    {#if matches.length > 15}
      <span class="pointer" on:click={() => page++} on:keypress={() => page++}
        >next...</span
      >
      {#if page != 0}
        <span class="pointer" on:click={() => page--} on:keypress={() => page--}
          >...previous</span
        >
      {/if}
    {/if}
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
          ><td class="center"
            >{matchinfo.time.slice(0, 16).replace("T", " ")}</td
          ></tr
        >
      </table>

      <table class="stats">
        <tr>
          <th /><th style="width: 20em;" /><th style="width: 3em;">K</th><th
            style="width: 3em;">A</th
          ><th style="width: 3em;">D</th>
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
          </tr>
        {/each}
      </table>
      <table id="rounds">
        <tr style="border-bottom: 1px solid black;">
          {#each matchinfo.rounds as round}
            <td
              >{#if round.notice == "halftime"}<div class="halftime">
                  |
                </div>{/if}{#if round.side == "CT"}<div class="roundct">
                  {@html getRoundLogo(round.notice)}
                </div>{/if}</td
            >
          {/each}
        </tr>
        <tr>
          {#each matchinfo.rounds as round}
            <td
              >{#if round.notice == "halftime"}<div class="halftime">
                  |
                </div>{/if}{#if round.side == "TERRORIST"}<div class="roundt">
                  {@html getRoundLogo(round.notice)}
                </div>{/if}</td
            >
          {/each}
        </tr>
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
  .pointer {
    cursor: pointer;
  }

  #rounds {
    margin-left: 100px;
    margin-top: 50px;
  }

  .roundct {
    background: linear-gradient(#434b50, #434b50, #b5d4ee33);
    color: #b5d4ee;
    padding: 0.3em;
  }

  .roundt {
    background: linear-gradient(#ead18a33, #434b50, #434b50);
    color: #ead18a;
    padding: 0.3em;
  }
</style>
