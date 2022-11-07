
  type Player = {
    name: string;
    kills: number;
    assists: number;
    deaths: number;
  };

type Round = {
  side: string;
  notice: string;
}
  export type MatchInfo = {
    filename: string;
    map: string;
    time: string;
    duration: number;
    players_ct: Player[];
    players_t: Player[];
    score_ct: number;
    score_t: number;
    rounds: Round[];
  };

  export type MatchBrief = {
    filename: string;
    map: string;
    time: string;
    score_ct: number;
    score_t: number;
  };
