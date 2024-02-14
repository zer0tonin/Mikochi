import { h } from "preact";
import Header from "../src/components/header";
// See: https://github.com/preactjs/enzyme-adapter-preact-pure
import { shallow } from "enzyme";

describe("Initial Test of the Header", () => {
  test("Header contains a 2 elem nav bar", () => {
    const context = shallow(<Header />);
    expect(context.find("nav").children().length).toBe(2);
  });
});
