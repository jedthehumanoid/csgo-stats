<script>
  export let match;
  export let selected = false;

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

<div on:click on:keypress class="match {selected ? 'selected' : ''}">
  <table>
    <tr>
      <td class="logo"
        ><img src="map-icons/map_icon_{match.map}.svg" alt="" /></td
      >
      <td
        ><span class="ct">{match.score_ct}</span> -
        <span class="terrorist">{match.score_t}</span><br />{getDate(
          match.filename
        )}</td
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
