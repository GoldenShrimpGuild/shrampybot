import type { StreamHistoryDatum } from '../../../model/utility/nosqldb/index';

// Data for sorting streams out
const gsgChannelLogin = "goldenshrimpguild";
const musicCategoryName = "Music";

export interface GSGStream extends StreamHistoryDatum {
  [key: string]: any,
  isEventStream?: boolean,
  isNormalStream?: boolean,
  raidTrain?: number,
}

export class GSGStreams<Key extends string, Stream extends GSGStream> {
  protected hiddenGSG: boolean
  protected gsgRecord: StreamHistoryDatum | null
  protected onlyMusic: boolean
  protected streams: Map<Key, Stream>
  protected nonMusicStreams: Map<Key, Stream>

  constructor(streamList?: Stream[], hideGSG?: boolean, onlyMusic?: boolean) {
    this.hiddenGSG = hideGSG || false
    this.onlyMusic = onlyMusic || false
    this.gsgRecord = null
    this.streams = new Map<Key, Stream>()
    this.nonMusicStreams = new Map<Key, Stream>()

    if (streamList) {
      streamList.forEach((stream: Stream) => {
        this.set(stream.user_login as Key, stream)
      })
    }
  }
  async reconcile(streamList: Stream[]) {
    const streamLogins = streamList.map(v => v.user_login)

    // Remove what's needed
    Array.from(this.streams.values())
      .filter((v) => !streamLogins.includes(v.user_login))
      .forEach((stream: Stream) => {
        this.streams.delete(stream.user_login as Key)
      })

    // Purges gsgRecord if no match to incoming list
    if (this.gsgRecord && !streamLogins.includes(this.gsgRecord.user_login)) {
      this.gsgRecord = null
    }

    // Purges non-music streams from list if no longer included overall
    Array.from(this.nonMusicStreams.values())
      .filter((v) => !streamLogins.includes(v.user_login))
      .forEach((stream: Stream) => {
        if (stream.game_name !== musicCategoryName) {
          // remove non-music streams from separate non-music map as well
          if (this.nonMusicStreams.has(stream.user_login as Key)) {
            this.nonMusicStreams.delete(stream.user_login as Key)
          }
        }
      })

    // Add what's needed
    streamList.forEach((stream: Stream) => {
      // Stow record (set is overridden so the gsgRecord and nonMusic logic are dealt with there)
      this.set(stream.user_login as Key, stream as Stream)
    })
  }
  setStream(stream: Stream) {
    if (!(stream as Stream || !stream.user_login)) {
      return this
    }
    return this.set(stream.user_login as Key, stream)
  }
  // Override existing set
  set(key: Key, value: Stream) {
    if (!key && !value) {
        // return without writing if no user_login or key is set
        return this
    }

    if (value.user_login === gsgChannelLogin) {
      this.gsgRecord = value
      if (this.hiddenGSG) {
        // return without writing
        return this
      }
    }

    if (value.game_name !== musicCategoryName) {
      this.nonMusicStreams.set(key, value)
      if (this.onlyMusic) {
        // return without writing
        return this
      }
    }

    // Stow record
    this.streams.set(key, value)
    return this
  }
  setHideGSG(hide: boolean) {
    if (hide) {
      Array.from(this.streams.values())
        .filter((v) => v.user_login === gsgChannelLogin)
        .forEach((stream: Stream) => {
          this.streams.delete(stream.user_login as Key)
        })
    } else if (this.gsgRecord) {
      this.streams.set(this.gsgRecord.user_login as Key, this.gsgRecord as Stream)
    }

    this.hiddenGSG = hide
    return this.hiddenGSG
  }
  // reshuffles streams to show only the Music category ones
  // returns map of non-music streams
  setOnlyMusic(onlyMusic: boolean) {
    if (onlyMusic) {
      const thisList: Stream[] = Array.from(this.streams.values())
        .filter((v: Stream) => v.game_name !== musicCategoryName)

      for (var i = 0; i < thisList.length; i++) {
        this.streams.delete(thisList[i].user_login as Key)
      }
    } else {
      const nmList = Array.from(this.nonMusicStreams.values())
      for (var i = 0; i < nmList.length; i++) {
        this.streams.set(nmList[i].user_login as Key, nmList[i] as Stream)
      }
    }

    this.onlyMusic = onlyMusic
    return this.nonMusicStreams
  }
  toggleHideGSG() {
    return this.setHideGSG(!this.hiddenGSG)
  }
  toggleOnlyMusic() {
    return this.setOnlyMusic(!this.onlyMusic)
  }
  isHidingGSG() {
    return this.hiddenGSG
  }
  keys() {
    return this.streams.keys()
  }
  values() {
    return this.streams.values()
  }
}