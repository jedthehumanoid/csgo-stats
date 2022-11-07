
  type Player = {
    name: string;
    kills: number;
    assists: number;
    deaths: number;
  };

  export type MatchInfo = {
    filename: string;
    map: string;
    time: string;
    duration: number;
    players_ct: Player[];
    players_t: Player[];
    score_ct: number;
    score_t: number;
  };

  export type MatchBrief = {
    filename: string;
    map: string;
    time: string;
    score_ct: number;
    score_t: number;
  };
