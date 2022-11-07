<script lang="ts">
  export let match: MatchBrief;
  export let selected = false;

  import type { MatchBrief } from "./csgo";
</script>

<div on:click on:keypress class="match {selected ? 'selected' : ''}">
  <table>
    <tr>
      <td class="logo"
        ><img src="map-icons/map_icon_{match.map}.svg" alt="" /></td
      >
      <td
        ><span class="ct">{match.score_ct}</span> -
        <span class="terrorist">{match.score_t}</span><br />{match.time
          .slice(0, 16)
          .replace("T", " ")}</td
      >
    </tr>
  </table>
</div>

<style>
  .ct {
    color: #b5d4ee;
  }

  .terrorist {
    color: #ead18a;
  }

  .match {
    cursor: pointer;
    border-bottom: 1px solid #000000;
  }

  .match:hover {
    background-color: #00000055;
  }

  .selected {
    background-color: #00000022;
    border-left: 4px solid #007ea7;
  }
  .logo {
    width: 75px;
  }

  img {
    width: 50px;
  }
</style>
