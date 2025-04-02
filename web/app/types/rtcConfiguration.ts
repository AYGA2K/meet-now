export interface RTCConfiguration {
  configuration: {
    offerToReceiveAudio: boolean;
    offerToReceiveVideo: boolean;
  };
  iceServers: Array<{ urls: string }>;
}
