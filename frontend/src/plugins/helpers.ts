export class helper {
  static pxToRem(px: number, base: number = 16) {
    return (1 / base) * px
  }
  static parsePixels(pxValue: string | number) {
    return parseInt(`${pxValue}`.replace(/px/, ""))
  }
  static ccToSlug(ccItem: string) {
    const nonSluggableRegex = /[^a-zA-Z0-9]/g
    return ccItem && ccItem[0].search(nonSluggableRegex) == -1 ? ccItem.replace(/([A-Z])/g, "-$1").toLowerCase() : ccItem
  }
}
