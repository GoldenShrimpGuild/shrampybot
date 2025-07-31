import type { StreamHistoryDatum } from '../../../model/utility/nosqldb/index';

// Data for sorting streams out
const gsgChannelLogin = "goldenshrimpguild";

export class Streams<Key extends string, Stream extends StreamHistoryDatum> extends Map {
  protected hiddenGSG: boolean
  protected gsgRecord: StreamHistoryDatum | null

  constructor(streamList?: Stream[], hideGSG?: boolean) {
    super()

    this.hiddenGSG = hideGSG || false
    this.gsgRecord = null

    if (streamList) {
      streamList.forEach((stream: Stream) => {
        this.set(stream.user_login as Key, stream)
      })
    }
  }
  reconcile(streamList: Stream[]) {
    const streamLogins = streamList.map(v => v.user_login)

    // Remove what's needed
    Array.from(this.values())
      .filter((v) => !streamLogins.includes(v.user_login))
      .forEach((stream: Stream) => {
        if (stream.user_login === gsgChannelLogin) {
          // clear gsgRecord when a matched stream is removed
          this.gsgRecord = null
        }
        stream.user_login
      })

    // Add what's needed
    streamList.forEach((stream: Stream) => {
      // Separately stow the GSG record and move on
      if (stream.user_login === gsgChannelLogin) {
        this.gsgRecord = stream
        return
      }
      // Stow record
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
    // Stow record
    super.set(key, value)
    return this
  }
  setHideGSG(hide: boolean) {
    if (hide) {
      Array.from(this.values())
        .filter((v) => v.user_login === gsgChannelLogin)
        .forEach((stream: Stream) => {
          this.delete(stream.user_login)
        })
    } else if (this.gsgRecord) {
      super.set(this.gsgRecord.user_login, this.gsgRecord)
    }

    this.hiddenGSG = hide
    return this.hiddenGSG
  }
  toggleHideGSG() {
    return this.setHideGSG(!this.hiddenGSG)
  }
  isHidingGSG() {
    return this.hiddenGSG
  }
}